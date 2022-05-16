package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type AircraftList struct {
	sync.RWMutex
	val string
}

func (l *AircraftList) Read() string {
	l.RLock()
	defer l.RUnlock()
	return l.val
}

func (l *AircraftList) Write(newVal string) {
	l.Lock()
	defer l.Unlock()
	l.val = newVal
}

var aircraftList = AircraftList{
	val: `
{
	"timestamp":0,
	"list":[]
}`,
}

func main() {
	var projectID = mustGetenv(env_projectID)

	go backgroundUpdateState(projectID, "aircraft-list", &aircraftList)

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/airspace/aircraft/list", getList)

	router.Run()
}

func getList(c *gin.Context) {

	c.String(http.StatusOK, aircraftList.Read())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("Environment variable not set: " + k)
	}
	return v
}

func backgroundUpdateState(projectId string, documentID string, state *AircraftList) {

	client := FirestoreInit(projectId)
	defer client.Close()

	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc(documentID).Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		snap, err := snapIter.Next()
		checkErr(err)

		if !snap.Exists() {
			panic("Document no longer exists.")
		}

		jsonData, err := json.Marshal(snap.Data())
		checkErr(err)

		state.Write(string(jsonData))
	}
}
