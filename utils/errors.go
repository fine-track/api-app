package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

func FailOnError(err error, msg string, c *gin.Context) {
	if err != nil {
		log.Panicf("%s | error: %s", msg, err.Error())
		if c != nil {
			c.Abort()
		}
	}
}
