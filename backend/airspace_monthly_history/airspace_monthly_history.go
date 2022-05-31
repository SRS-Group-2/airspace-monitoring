package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"github.com/gin-contrib/secure"
)

const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"
const logName = "MONTHLY_HISTORY_LOG"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

type LogType struct {
	Debug    *log.Logger
	Error    *log.Logger
	Critical *log.Logger
}

var Log = LogType{}

func main() {
	var projectID = mustGetenv(env_projectID)

	ctx := context.Background()
	loggerClient, err := logging.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer loggerClient.Close()

	Log.Debug = loggerClient.Logger(logName).StandardLogger(logging.Debug)
	Log.Error = loggerClient.Logger(logName).StandardLogger(logging.Error)
	Log.Critical = loggerClient.Logger(logName).StandardLogger(logging.Critical)

	Log.Debug.Print("Starting Monthly History Service.")
	defer Log.Debug.Println("Stopping Monthly History Service.")

	client := FirestoreInit(projectID)
	defer client.Close()

	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(func() gin.HandlerFunc {
        return func(c *gin.Context) {
            c.Writer.Header().Set("Cache-Control", "public, max-age=300")
        }
    }())

	router.Use(secure.New(secure.Config{
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

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
		var fromUtc string
		var toUtc string

		switch resolution {

		case "hour":
			fromUtc = from.UTC().Format("2006-01-02-15")
			toUtc = to.UTC().Format("2006-01-02-15")

			docIter = client.Collection("airspace/30d-history/1h-bucket").
				Where("startTime", ">=", fromUtc).
				Where("startTime", "<=", toUtc).
				OrderBy("startTime", firestore.Desc).
				Documents(ctx)

		case "day":
			fallthrough

		default:
			fromUtc = from.UTC().Format("2006-01-02")
			toUtc = to.UTC().Format("2006-01-02")

			docIter = client.Collection("airspace/30d-history/1d-bucket").
				Where("startTime", ">=", fromUtc).
				Where("startTime", "<=", toUtc).
				OrderBy("startTime", firestore.Desc).
				Documents(ctx)

			resolution = "day"
		}

		var jsonResult = make(map[string]interface{})
		var jsonHistory = make(map[string]interface{})

		// Iterate over documents and create response
		for i := 0; true; i++ {
			doc, err := docIter.Next()

			if err == iterator.Done {
				break
			}

			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				Log.Error.Println("Error iterating over history documents with ", resolution, " resolution: ", err)
				return
			}

			jsonHistory[doc.Ref.ID] = doc.Data()
		}

		jsonResult["from"] = fromUtc
		jsonResult["to"] = toUtc
		jsonResult["resolution"] = resolution
		jsonResult["history"] = jsonHistory

		c.JSON(http.StatusOK, jsonResult)
	})
	router.Run()
}

func checkErr(err error) {
	if err != nil {
		Log.Critical.Println("Critical error, panicking: ", err)
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
