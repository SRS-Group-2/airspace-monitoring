package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/jasonlvhit/gocron"
	"google.golang.org/api/iterator"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"
const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type HistoryValues struct {
	CO2       int64
	Distance  int64
	Timestamp string
}

type HistoryState struct {
	sync.RWMutex
	val HistoryValues
}

func (l *HistoryState) ReadCO2() int64 {
	l.RLock()
	defer l.RUnlock()
	return l.val.CO2
}

func (l *HistoryState) ReadDistance() int64 {
	l.RLock()
	defer l.RUnlock()
	return l.val.Distance
}

func (l *HistoryState) ReadTimestamp() string {
	l.RLock()
	defer l.RUnlock()
	return l.val.Timestamp
}

func (l *HistoryState) WriteCO2(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.val.CO2 = newVal
}

func (l *HistoryState) WriteDistance(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.val.Distance = newVal
}

func (l *HistoryState) WriteTimestamp(newVal string) {
	l.Lock()
	defer l.Unlock()
	l.val.Timestamp = newVal
}

func (l *HistoryState) IncCO2(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.val.CO2 += newVal
}

func (l *HistoryState) IncDistance(newVal int64) {
	l.Lock()
	defer l.Unlock()
	l.val.Distance += newVal
}

var oneHourState = &HistoryState{
	val: HistoryValues{
		CO2:       0,
		Distance:  0,
		Timestamp: "2006-01-01-00-00",
	},
}

var sixHoursState = &HistoryState{
	val: HistoryValues{
		CO2:       0,
		Distance:  0,
		Timestamp: "2006-01-01-00-00",
	},
}
var oneDayState = &HistoryState{
	val: HistoryValues{
		CO2:       0,
		Distance:  0,
		Timestamp: "2006-01-01-00-00",
	},
}

func main() {
	var credJson = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)

	client := FirestoreInit([]byte(credJson), projectID)
	defer client.Close()

	// Service just woke up, initialize the values of the states to the sum of what is in the DB

	lastUpdateTime := coldLoadFromDb(client)

	// Update the database entries for 1h, 6h 24h sums with initial values
	saveStateToDb(client)

	//decrease states values
	updateWithTimer(client)

	// Listen for new 5m data from db and update the state
	listenForDbUpdates(client, lastUpdateTime)
}

func coldLoadFromDb(client *firestore.Client) string {

	docIter := getAllFromDb(client)
	var lastUpdateTime = "2006-01-01-00-00"

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

		values := getValuesFromData(data)

		lastUpdateTime = incStateIfInRange(values, lastUpdateTime)
	}
	return lastUpdateTime
}

func listenForDbUpdates(client *firestore.Client, lastUpdateTime string) {
	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc("5m-history").Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		docSnap := waitForNewSnapshot(snapIter.Next())
		lastUpdateTime = documentUpdateHandler(docSnap, lastUpdateTime)
	}
}

func waitForNewSnapshot(docSnap *firestore.DocumentSnapshot, err error) *firestore.DocumentSnapshot {
	checkErr(err)

	if !docSnap.Exists() {
		panic("Document no longer exists.")
	}

	return docSnap
}

func documentUpdateHandler(docSnap *firestore.DocumentSnapshot, lastUpdateTime string) string {

	data := docSnap.Data()

	values := getValuesFromData(data)

	if isValidNewDocument(values.Timestamp, lastUpdateTime) {
		incStates(values)
		lastUpdateTime = values.Timestamp
	}
	return lastUpdateTime
}

func isValidNewDocument(timestamp string, lastUpdateTime string) bool {
	_, err := parseDate(timestamp)

	if err != nil {
		return false
	}

	if timestamp > lastUpdateTime {
		return true
	}

	return false
}

func getValuesFromData(data map[string]interface{}) HistoryValues {
	values := HistoryValues{}

	values.Timestamp = data["startTime"].(string)
	values.CO2 = data["CO2t"].(int64)
	values.Distance = data["distanceKm"].(int64)

	return values
}

func getAllFromDb(client *firestore.Client) *firestore.DocumentIterator {

	from := time.Now().UTC().AddDate(0, 0, -1) //from 24h
	fromUtc := from.UTC().Format("2006-01-02-15-04")
	toUtc := time.Now().UTC().Format("2006-01-02-15-04")

	ctx := context.Background()
	docIter := client.Collection("airspace/24h-history/5m-bucket").
		Where("startTime", ">=", fromUtc).
		Where("startTime", "<=", toUtc).
		OrderBy("startTime", firestore.Asc).
		Documents(ctx)

	return docIter
}

func incState(state *HistoryState, val HistoryValues) {
	state.IncCO2(val.CO2)
	state.IncDistance(val.Distance)
	state.WriteTimestamp(truncateToFiveMin(time.Now().UTC()).Format("2006-01-02-15-04"))
}

func incStates(val HistoryValues) {
	incState(oneHourState, val)
	incState(sixHoursState, val)
	incState(oneDayState, val)
}

func incStateIfInRange(values HistoryValues, lastUpdateTime string) string {
	referenceTime := truncateToFiveMin(time.Now().UTC().Truncate(time.Minute))

	if isAfter(values.Timestamp, referenceTime.Add(1*time.Hour)) {
		incState(oneHourState, values)
		lastUpdateTime = values.Timestamp
	}

	if isAfter(values.Timestamp, referenceTime.Add(6*time.Hour)) {
		incState(sixHoursState, values)
	}

	if isAfter(values.Timestamp, referenceTime.Add(24*time.Hour)) {
		incState(oneDayState, values)
	}
	return lastUpdateTime
}

func isAfter(timestamp string, threshold time.Time) bool {
	dateTime, err := parseDate(timestamp)
	if err != nil {
		fmt.Println("Error parsing date: ", timestamp, err)
		return false
	}

	//Add one minute here so that exactly matching timestamp are also included
	// 13:00 .After(13:00) -> false
	// 13:01 .After(13:00) -> true
	dateTime.Add(time.Minute)

	return dateTime.After(threshold)
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

func parseDate(timestamp string) (time.Time, error) {
	timeVal, err := time.Parse("2006-01-02-15-04", timestamp)

	if err != nil {
		return timeVal, err
	}

	if timeVal.Minute()%5 != 0 {
		return timeVal, errors.New("not a 5 minute time")
	}

	return timeVal, err
}

func cronjobHandler5min(client *firestore.Client) {
	// Since 5 minutes passed, subtract values that are now out of 1h,6h,24h interval from state
	now := truncateToFiveMin(time.Now().UTC())

	time24hAgo := now.AddDate(0, 0, -1).Format("2006-01-02-15-04")
	time6hAgo := now.Add(-6 * time.Hour).Format("2006-01-02-15-04")
	time1hAgo := now.Add(-1 * time.Hour).Format("2006-01-02-15-04")

	//remove value from 24h ago
	oneDayAgoDoc, err := getDocFromDB(client, time24hAgo)
	if err == nil {
		decreaseStateValues(oneDayAgoDoc.Data(), oneDayState)
	}

	//remove value from 6h ago
	sixHoursAgoDoc, err := getDocFromDB(client, time6hAgo)
	if err == nil {
		decreaseStateValues(sixHoursAgoDoc.Data(), sixHoursState)
	}

	//remove value from 1h ago
	oneHourAgoDoc, err := getDocFromDB(client, time1hAgo)
	if err == nil {
		decreaseStateValues(oneHourAgoDoc.Data(), oneHourState)
	}

	saveStateToDb(client)
}

func getDocFromDB(client *firestore.Client, documentID string) (*firestore.DocumentSnapshot, error) {
	ctx := context.Background()
	doc, err := client.Collection("airspace/24h-history/5m-bucket").Doc(documentID).Get(ctx)
	if err != nil {
		return doc, err
	}

	if !doc.Exists() {
		return doc, errors.New("document not found")
	}

	return doc, err
}

func decreaseStateValues(data map[string]interface{}, state *HistoryState) {
	values := HistoryValues{}

	values.CO2 = -data["CO2t"].(int64)
	values.Distance = -data["distanceKm"].(int64)

	incState(state, values)
}

func cronjobHandler1hour(client *firestore.Client) {
	// Write down 1h state to firestore 30d-history 1h-bucket
	ctx := context.Background()
	documentID := time.Now().Format("2006-01-02-15")
	client.Collection("airspace").Doc("30d-history").Collection("1h-bucket").Doc(documentID).Set(ctx, map[string]interface{}{
		"CO2t":       oneHourState.ReadCO2(),
		"distanceKm": oneHourState.ReadDistance(),
		"timestamp":  oneHourState.ReadTimestamp(),
	})
}

func cronjobHandler1day(client *firestore.Client) {

}

func truncateToFiveMin(dateTime time.Time) time.Time {
	trunc := dateTime.UTC().Truncate(time.Minute)
	rest := trunc.Minute() % 5
	dateTime.Add(-(time.Duration(rest) * time.Minute))
	return dateTime
}

func updateWithTimer(client *firestore.Client) {
	startTime5min := truncateToFiveMin(time.Now().UTC())
	startTime5min.Add(time.Minute * 5).Add(time.Minute)

	startTime1Hour := time.Now().UTC().Truncate(time.Hour).Add(time.Hour).Add(time.Minute * 2)
	startTime1Day := time.Now().UTC().Truncate(time.Hour*24).AddDate(0, 0, 1).Add(time.Minute * 3)

	gocron.Every(5).Minutes().From(&startTime5min).Do(cronjobHandler5min, client)
	gocron.Every(1).Hour().From(&startTime1Hour).Do(cronjobHandler1hour, client)
	gocron.Every(1).Day().From(&startTime1Day).Do(cronjobHandler1day, client)

}
