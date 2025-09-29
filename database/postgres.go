package database

import (
	"fmt"
	"os"

	"github.com/pubestpubest/pos-backend/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB(runEnv string) (err error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USERNAME"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	log.Info("[database]: Connected to database")

	db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Area{},
		&models.DiningTable{},
		&models.Category{},
		&models.MenuItem{},
		&models.Modifier{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderItemModifier{},
		&models.Payment{},
		&models.RolePermission{},
		&models.UserRole{},
	)
	log.Info("[database]: Migrated database")

	DB = db

	return err
}
