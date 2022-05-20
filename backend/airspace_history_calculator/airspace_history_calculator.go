package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"github.com/go-co-op/gocron"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"
const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const env_logName = "GOOGLE_LOG_NAME_HISTORY_CALCULATOR"
const env_cred = "GOOGLE_APPLICATION_CREDENTIALS"

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

func (l *HistoryState) Read() HistoryValues {
	l.RLock()
	defer l.RUnlock()
	result := HistoryValues{}
	result.CO2 = l.val.CO2
	result.Distance = l.val.Distance
	result.Timestamp = l.val.Timestamp
	return result
}

func (l *HistoryState) Add(addVal HistoryValues) {
	l.Lock()
	defer l.Unlock()
	l.val.CO2 += addVal.CO2
	l.val.Distance += addVal.Distance
	l.val.Timestamp = addVal.Timestamp
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

type LogType struct {
	Debug    *log.Logger
	Error    *log.Logger
	Critical *log.Logger
}

var Log = LogType{}

func main() {
	var projectID = mustGetenv(env_projectID)
	var logName = mustGetenv(env_logName)

	ctx := context.Background()
	loggerClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer loggerClient.Close()

	Log.Debug = loggerClient.Logger(logName).StandardLogger(logging.Debug)
	Log.Error = loggerClient.Logger(logName).StandardLogger(logging.Error)
	Log.Critical = loggerClient.Logger(logName).StandardLogger(logging.Critical)

	Log.Debug.Print("Starting History Calculator Service.")
	defer Log.Debug.Println("Stopping History Calculator Service.")

	client := FirestoreInit(projectID)
	defer client.Close()

	// Service just woke up, initialize the values of the states to the sum of what is in the DB
	lastUpdateTime := coldLoadFromDb(client)

	// Update the database entries for 1h, 6h 24h sums with initial values
	saveStateToDb(client)

	//decrease states values, update 1h, 1d, values, cleanup db.
	startCronJobs(client)

	// Listen for new 5m data from db and update the state
	listenForDbUpdates(client, lastUpdateTime)

}

func coldLoadFromDb(client *firestore.Client) string {
	Log.Debug.Println("Initializing from database values.")

	docIter := getAllFromDb(client)
	var lastUpdateTime = "2006-01-01-00-00"

	// Iterate over the 5m values of the last 24h and add them to the states if relevant
	for {
		doc, err := docIter.Next()

		if err == iterator.Done {

			break
		}

		if err != nil {
			code := status.Code(err)
			if code == codes.PermissionDenied || code == codes.Unauthenticated {
				panic(err)
			} else {
				fmt.Println("Error calculating initial state, iterating documents:", err)
				continue
			}
		}

		data := doc.Data()

		values := getValuesFromData(data)

		lastUpdateTime = incStateIfInRange(values, lastUpdateTime)

	}
	return lastUpdateTime
}

func getAllFromDb(client *firestore.Client) *firestore.DocumentIterator {

	from := time.Now().UTC().AddDate(0, 0, -1) //from 24h
	fromUtc := truncTo5Min(from.UTC()).Format("2006-01-02-15-04")

	ctx := context.Background()
	docIter := client.Collection("airspace/24h-history/5m-bucket").
		Where("startTime", ">=", fromUtc).
		OrderBy("startTime", firestore.Asc).
		Documents(ctx)

	return docIter
}

func listenForDbUpdates(client *firestore.Client, lastUpdateTime string) {

	ctx := context.Background()
	snapIter := client.Collection("airspace").Doc("5m-history").Snapshots(ctx)

	for {
		// Wait for new snapshots of the document
		docSnap := waitForNewSnapshot(snapIter.Next())
		lastUpdateTime = documentUpdateHandler(client, docSnap, lastUpdateTime)

	}
}

func waitForNewSnapshot(docSnap *firestore.DocumentSnapshot, err error) *firestore.DocumentSnapshot {

	checkErr(err)

	if !docSnap.Exists() {
		panic("Document no longer exists.")
	}

	return docSnap
}

func documentUpdateHandler(client *firestore.Client, docSnap *firestore.DocumentSnapshot, lastUpdateTime string) string {

	data := docSnap.Data()

	values := getValuesFromData(data)

	if isValidNewDocument(values.Timestamp, lastUpdateTime) {
		lastUpdateTime = values.Timestamp
		incStates(values)
		saveStateToDb(client)
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

func incStates(val HistoryValues) {

	oneHourState.Add(val)
	sixHoursState.Add(val)
	oneDayState.Add(val)
}

func incStateIfInRange(values HistoryValues, lastUpdateTime string) string {

	referenceTime := truncTo5Min(time.Now().UTC().Truncate(time.Minute))

	if isAfter(values.Timestamp, referenceTime.Add(-1*time.Hour)) {
		oneHourState.Add(values)
	}

	if isAfter(values.Timestamp, referenceTime.Add(-6*time.Hour)) {
		sixHoursState.Add(values)
	}

	if isAfter(values.Timestamp, referenceTime.Add(-24*time.Hour)) {
		oneDayState.Add(values)
		lastUpdateTime = values.Timestamp
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

	now := truncTo5Min(time.Now().UTC())

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
	values.Timestamp = truncTo5Min(time.Now().UTC()).Format("2006-01-02-15-04")

	state.Add(values)
}

func cronjobHandler1hour(client *firestore.Client) {

	// Write down 1h state to firestore 30d-history 1h-bucket
	ctx := context.Background()
	documentID := time.Now().UTC().Format("2006-01-02-15")
	startTime := time.Now().UTC().Add(-time.Hour).Format("2006-01-02-15")

	client.Collection("airspace").Doc("30d-history").Collection("1h-bucket").Doc(documentID).Set(ctx, map[string]interface{}{
		"CO2t":       oneHourState.ReadCO2(),
		"distanceKm": oneHourState.ReadDistance(),
		"startTime":  startTime,
	})
}

func cronjobHandler1day(client *firestore.Client) {

	// Write down 1d state to firestore 30d-history 1d-bucket
	ctx := context.Background()

	// We are recording the data of the previous day
	documentID := time.Now().UTC().AddDate(0, 0, -1).Format("2006-01-02")

	client.Collection("airspace").Doc("30d-history").Collection("1d-bucket").Doc(documentID).Set(ctx, map[string]interface{}{
		"CO2t":       oneHourState.ReadCO2(),
		"distanceKm": oneHourState.ReadDistance(),
		"startTime":  documentID,
	})

	// Cleanup some old documents from the collection
	go cleanupDb(client)
}

func cleanupDb(client *firestore.Client) {

	// Start deleting documents older than 25h ago in 24h-history/5m-bucket
	threshold5min := time.Now().UTC().AddDate(0, 0, -1).Add(-time.Hour)
	deleteOlderThanFrom(client, threshold5min, "24h-history/5m-bucket")

	// Start deleting documents older than 31d ago in 30d-history/1h-bucket
	thresholdHours := time.Now().UTC().AddDate(0, 0, -31)
	deleteOlderThanFrom(client, thresholdHours, "30d-history/1h-bucket")

	// Start deleting documents older than 32d ago in 30dh-history/1d-bucket
	thresholdDays := time.Now().UTC().AddDate(0, 0, -32)
	deleteOlderThanFrom(client, thresholdDays, "30d-history/1d-bucket")
}

func deleteOlderThanFrom(client *firestore.Client, threshold time.Time, collectionPath string) {

	ctx := context.Background()
	docIter := client.Collection("airspace/"+collectionPath).
		Where("startTime", "<=", threshold).
		Documents(ctx)

	batchOp := client.Batch()

	for {
		doc, err := docIter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			continue
		}

		batchOp.Delete(doc.Ref)
	}

	_, err := batchOp.Commit(ctx)

	if err != nil {
		fmt.Println("Error deleting documents: ", err)
	}
}

func truncTo5Min(dateTime time.Time) time.Time {

	trunc := dateTime.UTC().Truncate(5 * time.Minute)

	return trunc
}

func startCronJobs(client *firestore.Client) {

	startTime5min := time.Now().UTC().Truncate(5 * time.Minute)
	startTime5min.Add(time.Minute * 5).Add(time.Minute)

	startTime1Hour := time.Now().UTC().Truncate(time.Hour).Add(time.Hour).Add(time.Minute * 2)
	startTime1Day := time.Now().UTC().Truncate(time.Hour*24).AddDate(0, 0, 1).Add(time.Minute * 3)

	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(5).Minutes().StartAt(startTime5min).Do(cronjobHandler5min, client)
	scheduler.Every(1).Hour().StartAt(startTime1Hour).Do(cronjobHandler1hour, client)
	scheduler.Every(1).Day().StartAt(startTime1Day).Do(cronjobHandler1day, client)

	scheduler.StartAsync()
}
