package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pubestpubest/pos-backend/database"
	"github.com/pubestpubest/pos-backend/middlewares"
	"github.com/pubestpubest/pos-backend/routes"
	"github.com/pubestpubest/pos-backend/seed"
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

	deployEnv := os.Getenv("DEPLOY_ENV")
	if deployEnv == "" {
		deployEnv = "local"
	}

	if deployEnv == "local" {
		if err := godotenv.Load("configs/.env"); err != nil {
			log.Fatal("[init]: Error loading .env file: ", err)
		}
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

	if err := database.ConnectMinio(); err != nil {
		log.Fatal("[init]: Connect MinIO error: ", err.Error())
	}

	isSeed := os.Getenv("SEED_DB")
	if isSeed == "" {
		isSeed = "false"
	}
	log.Info("[init]: Seed database: ", isSeed)
	if isSeed == "true" {
		seedRunner := seed.Runner{
			DB:  database.DB,
			Env: runEnv,
		}
		if err := seedRunner.Run(context.Background()); err != nil {
			log.Fatal("[init]: Seed database error: ", err.Error())
		}
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
	routes.AuthRoutes(v1)
	routes.CategoryRoutes(v1)
	routes.AreaRoutes(v1)
	routes.ModifierRoutes(v1)
	routes.OrderRoutes(v1)
	routes.PaymentRoutes(v1)
	routes.RoleRoutes(v1)
	routes.PermissionRoutes(v1)
	routes.UserRoutes(v1)
	routes.MenuItemRoutes(v1)
	routes.TableRoutes(v1)
	app.Run(":8080")
}
