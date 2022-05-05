package main

import (
	"context"
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

func main() {

	var credString = mustGetenv(env_credJson)
	var projectID = mustGetenv(env_projectID)

	//DB
	client := FirestoreInit([]byte(credString), projectID)

	router := gin.New()
	router.SetTrustedProxies(nil)

	// eg. airspace/history?from=UNIXTS&to=UNIXTS&resolution=
	// resolution: hour, day
	router.GET("airspace/history", func(c *gin.Context) {
		toDefault := time.Now().UTC()
		fromDefault := time.Now().UTC().AddDate(0, 0, -30)
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
			fromUtc := from.UTC().Format("2006-01-02-15")
			toUtc := to.UTC().Format("2006-01-02-15")

			docIter = client.Collection("airspace/30d-history/1h-bucket").
				Where("startTime", ">=", fromUtc).
				Where("startTime", "<=", toUtc).
				OrderBy("startTime", firestore.Desc).
				Documents(ctx)

		case "day":
			fallthrough

		default:
			fromUtc := from.UTC().Format("2006-01-02")
			toUtc := to.UTC().Format("2006-01-02")

			docIter = client.Collection("airspace/30d-history/1d-bucket").
				Where("startTime", ">=", fromUtc).
				Where("startTime", "<=", toUtc).
				OrderBy("startTime", firestore.Desc).
				Documents(ctx)
		}

		var json = make(map[string]interface{})

		// Iterate over documents and create response
		for i := 0; true; i++ {
			doc, err := docIter.Next()

			if err == iterator.Done {
				break
			}

			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				return
			}

			json[doc.Ref.ID] = doc.Data()
		}
		c.IndentedJSON(http.StatusOK, json)
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

	return time.Unix(result, 0).UTC()
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic("Environment variable not set: " + k)
	}
	return v
}
