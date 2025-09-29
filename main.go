package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pubestpubest/go-clean-arch-template/database"
	"github.com/pubestpubest/go-clean-arch-template/middlewares"
	"github.com/pubestpubest/go-clean-arch-template/routes"
	log "github.com/sirupsen/logrus"
)

var runEnv string

func init() {
	fmt.Println("Hello, World from init()")

	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetLevel(log.InfoLevel)

	runEnv = os.Getenv("RUN_ENV")
	if runEnv == "" {
		runEnv = "development"
	}

	if err := godotenv.Load("configs/.env"); err != nil {
		log.Fatal("[init]: Error loading .env file: ", err)
	}

	if runEnv == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Info("[init]: Run environment: ", runEnv)

	if err := database.ConnectDB(runEnv); err != nil {
		log.Fatal("[init]: Connect database PG error: ", err.Error())
	}
}

func main() {
	fmt.Println("Hello, World from main()")

	app := gin.Default()

	app.Use(middlewares.CORSMiddleware())

	app.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	app.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "not found",
		})
	})
	app.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"status": "method not allowed",
		})
	})

	v1 := app.Group("/v1")
	routes.UserRoutes(v1)

	app.Run(":8080")
}
