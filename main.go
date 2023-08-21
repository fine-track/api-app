package main

import (
	"log"
	"os"

	"github.com/fine-track/api-app/journals"
	"github.com/fine-track/api-app/ledgers"
	"github.com/fine-track/api-app/middlewares"
	"github.com/fine-track/api-app/profiles"
	"github.com/fine-track/api-app/utils"
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
	r.GET("/health", GetHealthHandler)

	v1 := r.Group("/v1")
	journalsGrp := v1.Group("/journals")
	{
		journalsGrp.POST("/", journals.CreateJournalHandler)
		journalsGrp.GET("/", journals.GetJournalsHandler)
		journalsGrp.PUT("/:id", journals.UpdateJournalHandler)
	}
	ledgersGrp := v1.Group("/ledgers")
	{
		ledgersGrp.GET("/", ledgers.GetLedgersHandler)
		ledgersGrp.PUT("/:id", ledgers.UpdateLedgerHandler)
		ledgersGrp.DELETE("/:id", ledgers.ResetLedgerHandler)
	}
	profileGrp := v1.Group("/profile")
	{
		profileGrp.GET("/:id", profiles.GetFullProfile)
		profileGrp.PUT("/:id", profiles.UpdateProfile)
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("unable to load .env file: %v", err.Error())
	}
	PORT := os.Getenv("PORT")

	journalsConn := journals.PrepareConn()
	// ledgersConn := ledgers.PrepareConn()

	defer func() {
		journalsConn.Close()
		// ledgersConn.Close()
	}()

	r := setupRouter()
	registerRoutes(r)

	err := r.Run(":" + PORT)
	utils.FailOnError(err, "unable to start server", nil)
}
