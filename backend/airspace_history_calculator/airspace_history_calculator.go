package main

import (
	"context"
	"encoding/json"
	"os"
	"sync"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"
const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type HistoryState struct {
	sync.RWMutex
	co2       int64
	distance  int64
	timestamp string
}

func (l *HistoryState) ReadCO2() int64 {
	l.RLock()
	defer l.RUnlock()
	return l.co2
}

func (l *HistoryState) ReadDistance() int64 {
	l.RLock()
	defer l.RUnlock()
	return l.distance
}

func (l *HistoryState) ReadTimestamp() string {
	l.RLock()
	defer l.RUnlock()
	return l.timestamp
}

func (l *HistoryState) WriteCO2(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.co2 = newVal
}

func (l *HistoryState) WriteDistance(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.distance = newVal
}

func (l *HistoryState) WriteTimestamp(newVal string) {
	l.Lock()
	defer l.Unlock()
	l.timestamp = newVal
}

var oneHourState = HistoryState{
	co2:       0,
	distance:  0,
	timestamp: "0000-00-00-00-00",
}

var sixHoursState = HistoryState{
	co2:       0,
	distance:  0,
	timestamp: "0000-00-00-00-00",
}
var oneDayState = HistoryState{
	co2:       0,
	distance:  0,
	timestamp: "0000-00-00-00-00",
}

func main() {
	var credJson = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)

	go backgroundUpdateState(credJson, projectID, "1h-history", &oneHourState)
}

func backgroundUpdateState(credString string, projectId string, documentID string, historyState *HistoryState) {

	client := FirestoreInit([]byte(credString), projectId)
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

		historyState.Write(string(jsonData))
	}
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
