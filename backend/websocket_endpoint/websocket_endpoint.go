package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const env_websocketEndpoint = "WEBSOCKET_ENDPOINT"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

var websocketEndpoint = mustGetenv(env_websocketEndpoint)

func main() {
	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/endpoint/position/url", getWebsocketEndpoint)

	router.Run()
}

func getWebsocketEndpoint(c *gin.Context) {
	c.String(http.StatusOK, websocketEndpoint)
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
