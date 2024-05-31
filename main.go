package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", HomeHandler)
	r.GET("/login", LoginHandler)
	r.GET("/oauth2/callback", CallbackHandler)
	r.GET("/profile", ProfileHandler)

	r.Run(":8181")
}
