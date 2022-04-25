package main

import (
	"database/sql"

	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Aircraft struct {
	Icao24           string `json:"icao24"`
	Manufacturername string `json:"manufacturername"`
	Model            string `json:"model"`
	Registration     string `json:"registration"`
	SerialNumber     string `json:"serialnumber"`
}

type Param struct {
	Icao24 string `form:"icao24" json:"icao24"`
}

const database_file = "./resources/aircraft_info.db"

var db *sql.DB
var queryStmt *sql.Stmt

func main() {

	database, err := sql.Open("sqlite3", database_file)
	checkErr(err)
	db = database

	stmt, err := db.Prepare(`SELECT 
					 icao24,
					 manufacturername, 
					 model, 
					 registration, 
					 serialnumber 
					 FROM aircraft_info 
					 WHERE icao24 =?;`)
	checkErr(err)

	queryStmt = stmt

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/aircraft/:icao24/info", getInfo)

	router.Run("localhost:8080")
}

func getInfo(c *gin.Context) {
	var icao24 = c.Param("icao24")
	if len(icao24) == 6 {
	} else {
		c.IndentedJSON(http.StatusNotAcceptable, gin.H{"message": "invalid icao24"})
		return
	}

	var rows, err = queryStmt.Query(icao24)
	checkErr(err)

	var aircraft Aircraft

	for rows.Next() {
		err = rows.Scan(&aircraft.Icao24,
			&aircraft.Manufacturername,
			&aircraft.Model,
			&aircraft.Registration,
			&aircraft.SerialNumber)
		checkErr(err)
	}

	if aircraft.Icao24 == "" {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "icao24 not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, aircraft)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
