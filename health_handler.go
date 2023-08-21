package main

import (
	"github.com/fine-track/api-app/utils"
	"github.com/gin-gonic/gin"
)

func GetHealthHandler(c *gin.Context) {
	res := utils.HTTPResponse{
		Message: "Health is OK",
	}
	user, _ := c.Get("user")
	res.Data = user
	res.Ok(c)
}
