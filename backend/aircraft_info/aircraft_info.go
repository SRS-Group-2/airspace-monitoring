package main

import (
	"net/http"

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

func main() {

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
		if len(icao24) != 6 {
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
		panic(err)
	}
}
