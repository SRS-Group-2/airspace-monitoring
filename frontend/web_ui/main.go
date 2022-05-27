package main

import (
	"os"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
)

const env_port = "PORT"
const env_ginmode = "GIN_MODE"

func main() {
	router := gin.New()
	router.SetTrustedProxies(nil)

	router.Use(func() gin.HandlerFunc {
        return func(c *gin.Context) {
            c.Writer.Header().Set("Cache-Control", "public, max-age=604800, immutable")
        }
    }())

	router.Use(secure.New(secure.Config{
		// AllowedHosts:          []string{"example.com", "ssl.example.com"},
		// SSLRedirect:           true,
		// SSLHost:               "ssl.example.com",
		STSSeconds:            315360000,
		STSIncludeSubdomains:  true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self' 'unsafe-inline' wss:",
		// IENoOpen:              true,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		// SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
	}))

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
