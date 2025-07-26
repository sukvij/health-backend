package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	chatservice "github.com/sukvij/health-checker/backend/chat-service"
	entiremedicalreport "github.com/sukvij/health-checker/backend/entire-medical-report"
	healthreportservice "github.com/sukvij/health-checker/backend/health-report-service"
	"github.com/sukvij/health-checker/backend/healthfers/database"
	userservice "github.com/sukvij/health-checker/backend/user-service"
)

func main() {
	db, _ := database.Connection()
	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // ALLOW ALL ORIGINS FOR DEVELOPMENT. BE MORE SPECIFIC IN PRODUCTION.
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // Include Authorization if you use it
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	userservice.UserServiceRoute(app, db)
	healthreportservice.HealthReportServiceController(app, db)
	chatservice.HistroyRoute(app, db)
	entiremedicalreport.EntireMedicalReportRoute(app, db)
	app.Run(":8080")
}
