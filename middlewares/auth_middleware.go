package middlewares

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/fine-track/api-app/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	AUTH_APP := os.Getenv("AUTH_APP")

	res := utils.HTTPResponse{}
	authToken := c.Request.Header.Get("authorization")
	if authToken == "" {
		res.Message = "No authorization token found"
		res.Unauthorized(c)
		c.Abort()
		return
	}

	splitted := strings.Split(authToken, " ")
	if splitted[0] != "Bearer" || len(splitted) < 2 || len(splitted) > 2 {
		res.Message = "Invalid authorization token format"
		res.Unauthorized(c)
		c.Abort()
		return
	}

	client := http.DefaultClient
	req, err := http.NewRequest(http.MethodGet, AUTH_APP+"/profile", nil)
	if err != nil {
		res.Message = err.Error()
		res.Unauthorized(c)
		c.Abort()
		return
	}
	req.Header.Set("Authorization", authToken)
	resp, err := client.Do(req)
	if err != nil {
		res.Message = err.Error()
		res.Unauthorized(c)
		c.Abort()
		return
	}
	defer resp.Body.Close()
	respData := utils.HTTPResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		res.Message = err.Error()
		res.Unauthorized(c)
		c.Abort()
		return
	}
	if resp.StatusCode != http.StatusOK {
		respData.Unauthorized(c)
		c.Abort()
		return
	}
	user := respData.Data
	c.Set("user", &user)
	c.Next()
}
