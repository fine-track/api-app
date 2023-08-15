package main

import (
	"log"
	"os"

	"github.com/fine-track/api-app/handlers"
	"github.com/fine-track/api-app/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://finetrack.sifatul.com", "http://localhost:3000"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Access-Control-Request-Headers",
		"Access-Control-Request-Method",
		"Authorization",
		"Connection",
		"Content-Type",
		"Cookie",
		"Date",
		"If-Modified-Since",
		"If-None-Match",
		"Origin",
		"Referrer",
		"User-Agent",
		"X-Requested-With",
	}
	r.Use(cors.New(corsConfig))
	r.Use(middlewares.AuthMiddleware)

	return r
}

func registerRoutes(r *gin.Engine) {
	r.GET("/health", handlers.GetHealthHandler)
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load .env file: %v", err.Error())
	}
	r := setupRouter()
	registerRoutes(r)
	PORT := os.Getenv("PORT")
	if err := r.Run(":" + PORT); err != nil {
		log.Fatalf("unable to run the server: %v", err.Error())
	}
}
