package main

import (
	"os"
	"github.com/gin-gonic/gin"
)

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

func main() {
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Static("/", "./statics")

	router.Run()
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