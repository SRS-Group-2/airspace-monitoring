package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/secure"
)

const env_websocketEndpoint = "WEBSOCKET_ENDPOINT"

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

var websocketEndpoint = mustGetenv(env_websocketEndpoint)

func main() {
	router := gin.New()

	router.SetTrustedProxies(nil)
	router.GET("/endpoints/position/url", getWebsocketEndpoint)

	router.Use(secure.New(secure.Config{
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	router.Run()
}

func getWebsocketEndpoint(c *gin.Context) {
	c.Header("Content-Type", "text/plain")
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
