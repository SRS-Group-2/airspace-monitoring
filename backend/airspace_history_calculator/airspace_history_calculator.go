package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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

func (l *HistoryState) IncCO2(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.co2 += newVal
}

func (l *HistoryState) IncDistance(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.distance += newVal
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

	client := FirestoreInit([]byte(credJson), projectID)
	defer client.Close()

	ctx := context.Background()

	// Service just woke up, initialize the values of the states to the sum of what is in the DB

	from := time.Now().UTC().AddDate(0, 0, -1)
	fromUtc := from.UTC().Format("2006-01-02-15")
	toUtc := time.Now().UTC().Format("2006-01-02-15")

	docIter := client.Collection("airspace/24h-history/5m-bucket").
		Where("startTime", ">=", fromUtc).
		Where("startTime", "<=", toUtc).
		OrderBy("startTime", firestore.Desc).
		Documents(ctx)

	// time.Now() but floored to nearest 5 minutes for comparisons
	begin, _ := time.Parse("2006-01-02-15", time.Now().UTC().Format("2006-01-02-15"))

	// Iterate over the 5m values of the last 24h and add them to the states if relevant
	for {
		doc, err := docIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("Error calculating initial state, iterating documents:", err)
			continue
		}

		data := doc.Data()
		timestamp := data["startTime"].(string)
		if inRange(timestamp, begin, 1*time.Hour) {
			oneHourState.IncCO2(data["CO2t"].(int64))
			oneHourState.IncDistance(data["distanceKm"].(int64))
			oneHourState.WriteTimestamp(timestamp)
		}
		if inRange(timestamp, begin, 6*time.Hour) {
			sixHoursState.IncCO2(data["CO2t"].(int64))
			sixHoursState.IncDistance(data["distanceKm"].(int64))
			sixHoursState.WriteTimestamp(timestamp)
		}
		if inRange(timestamp, begin, 24*time.Hour) {
			oneDayState.IncCO2(data["CO2t"].(int64))
			oneDayState.IncDistance(data["distanceKm"].(int64))
			oneDayState.WriteTimestamp(timestamp)
		}
	}

	// Update the database entries for 1h, 6h 24h sums with initial values
	saveStateToDb(client)

	// Now we start listenign to changes on the 5m-history document
	snapIter := client.Collection("airspace").Doc("5m-history").Snapshots(ctx)

	fistValue, err := snapIter.Next()
	checkErr(err)

	if !fistValue.Exists() {
		panic("5m update document no longer exists.")
	}

	data := fistValue.Data()
	timestamp := data["startTime"].(string)
	startTime, err := time.Parse("2006-01-02-15", timestamp)
	if err == nil {
		if startTime.After(begin) {
			updateState(client, data, startTime)
		}
	}

	for {
		// Wait for new snapshots of the document
		new5mins, err := snapIter.Next()
		checkErr(err)

		if !new5mins.Exists() {
			panic("Document no longer exists.")
		}

		new5Values := new5mins.Data()

	}
}

func updateState(client *firestore.Client, addData map[string]interface{}, startTime time.Time) {

}

func saveStateToDb(client *firestore.Client) {
	ctx := context.Background()

	client.Collection("airspace").Doc("24h-history").Update(ctx, []firestore.Update{
		{Path: "CO2t", Value: oneDayState.ReadCO2()},
		{Path: "distanceKm", Value: oneDayState.ReadDistance()},
		{Path: "timestamp", Value: oneDayState.ReadTimestamp()},
	})

	client.Collection("airspace").Doc("6h-history").Update(ctx, []firestore.Update{
		{Path: "CO2t", Value: sixHoursState.ReadCO2()},
		{Path: "distanceKm", Value: sixHoursState.ReadDistance()},
		{Path: "timestamp", Value: sixHoursState.ReadTimestamp()},
	})

	client.Collection("airspace").Doc("1h-history").Update(ctx, []firestore.Update{
		{Path: "CO2t", Value: oneHourState.ReadCO2()},
		{Path: "distanceKm", Value: oneHourState.ReadDistance()},
		{Path: "timestamp", Value: oneHourState.ReadTimestamp()},
	})
}

func backgroundUpdateState(credString string, projectId string, historyState *HistoryState) {

	client := FirestoreInit([]byte(credString), projectId)
	defer client.Close()

	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc("5m-history").Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		snap, err := snapIter.Next()
		checkErr(err)

		if !snap.Exists() {
			panic("Document no longer exists.")
		}

		snap.Data()

		historyState.Write()
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

func inRange(timestamp string, now time.Time, duration time.Duration) bool {
	sample, err := time.Parse("2006-01-02-15", timestamp)
	if err != nil {
		fmt.Println("Error parsing date: ", timestamp)
		return false
	}

	if sample.After(begin) {
		return true
	}
	return false

}
