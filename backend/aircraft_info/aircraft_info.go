package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"unicode"

	"cloud.google.com/go/logging"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/gin-gonic/gin"
)

type Aircraft struct {
	Icao24       string `json:"icao24"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	Registration string `json:"registration"`
	SerialNumber string `json:"serialnumber"`
}

type Param struct {
	Icao24 string `form:"icao24" json:"icao24"`
}

const max_db_connections = 25
const database_path = "file:./resources/aircraft_info.db"
const database_config = "?immutable=1&mode=ro"

const logName = "AIRCRAFT_INFO_LOG"
const env_cred = "GOOGLE_APPLICATION_CREDENTIALS"
const env_projectID = "GOOGLE_CLOUD_PROJECT_ID"

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

	Log.Debug.Print("Starting Aircraft Info Service.")
	defer Log.Debug.Println("Stopping Aircraft Info Service.")

	// Create multiple database connections to allow concurrent queries
	connections := make(chan *sqlite3.Stmt, max_db_connections)

	for i := 0; i < max_db_connections; i++ {
		conn, err := sqlite3.Open(database_path + database_config)
		checkErr(err)

		stmt, err := conn.Prepare(`SELECT
			icao24,
			manufacturername,
			model,
			registration,
			serialnumber
			FROM aircraft_info
			WHERE icao24 =?;`)
		checkErr(err)

		defer func() {
			stmt.Close()
			conn.Close()
		}()

		connections <- stmt
	}

	checkoutStmt := func() *sqlite3.Stmt {
		return <-connections
	}

	checkinStmt := func(c *sqlite3.Stmt) {
		connections <- c
	}

	router := gin.New()

	router.SetTrustedProxies(nil)

	router.GET("/airspace/aircraft/:icao24/info", func(c *gin.Context) {
		var icao24 = c.Param("icao24")

		if len(icao24) != 6 && hasSymbol(icao24) {
			c.String(http.StatusNotAcceptable, "Invalid icao24")
			return
		}

		stmt := checkoutStmt()
		defer checkinStmt(stmt)

		stmt.Bind(icao24)
		hasRow, err := stmt.Step()
		defer stmt.Reset()

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			Log.Error.Println("Error querying database for icao24=", icao24, " err: ", err)
			return
		}

		if !hasRow {
			c.String(http.StatusNotFound, "icao24 not found")
			return
		}

		var aircraft Aircraft

		stmt.Scan(&aircraft.Icao24,
			&aircraft.Manufacturer,
			&aircraft.Model,
			&aircraft.Registration,
			&aircraft.SerialNumber)

		c.JSON(http.StatusOK, aircraft)
	})

	router.Run()
}

func checkErr(err error) {
	if err != nil {
		Log.Critical.Println("Critical error, panicking: ", err)
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

func hasSymbol(str string) bool {
	for _, letter := range str {
		if unicode.IsSymbol(letter) {
			return true
		}
	}
	return false
}
