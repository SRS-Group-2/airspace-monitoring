package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

const env_credJson = "GOOGLE_APPLICATION_CREDENTIALS"

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type TimeResolution struct {
	oneHour string
	oneDay  string
}

var timeResolution = TimeResolution{
	oneHour: "hour",
	oneDay:  "day",
}

var maxDuration = time.Hour * 24 * 30

func main() {

	var credFile = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)

	//DB
	client := FirestoreInit(credFile, projectID)

	router := gin.New()
	router.SetTrustedProxies(nil)

	// eg. airspace/history?from=UNIXTS&to=UNIXTS&resolution=
	// resolution: hour, day
	//?type=co2
	router.GET("airspace/history", func(c *gin.Context) {
		toDefault := time.Now()
		fromDefault := time.Now().Add(-maxDuration)
		resolutionDefault := "day"

		fromStr := c.Query("from")
		from := parseDate(fromStr, fromDefault)

		toStr := c.Query("to")
		to := parseDate(toStr, toDefault)

		resolution := c.DefaultQuery("resolution", resolutionDefault)

		ctx := context.Background()

		var docIter *firestore.DocumentIterator

		switch resolution {
		case "hour":
			fromUtc := from.UTC().Format("YYYY-MM-DD-HH")
			toUtc := to.UTC().Format("YYYY-MM-DD-HH")
			docIter = client.Collection("airspace/30-history/1h-history").Where("startTime", ">=", fromUtc).Where("startTime", "<=", toUtc).Documents(ctx)
		case "day":
			fallthrough
		default:
			fromUtc := from.UTC().Format("YYYY-MM-DD")
			toUtc := to.UTC().Format("YYYY-MM-DD")
			docIter = client.Collection("airspace/30-history/1d-history").Where("startTime", ">=", fromUtc).Where("startTime", "<=", toUtc).Documents(ctx)
		}

		for {
			doc, err := docIter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				c.String(http.StatusInternalServerError, "There was an unexpected error while iterating")
				return
			}
			fmt.Println(doc.Data())
		}
		c.IndentedJSON(http.StatusOK, "hello")
	})

	router.Run()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parseDate(param string, defaultVal time.Time) time.Time {
	if param == "" {
		return defaultVal
	}

	result, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return defaultVal
	}

	return time.Unix(result, 0)
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("Environment variable not set: " + k)
	}
	return v
}
