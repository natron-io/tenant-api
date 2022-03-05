package database

import (
	"fmt"

	"github.com/natron-io/tenant-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConn      *gorm.DB
	err         error
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_SSLMODE  string
)

func InitDB() error {
	// Connect to the database
	dbUri := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE)
	DBConn, err = gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	err = DBConn.AutoMigrate(
		&models.Tenant{},
		&models.CPUCost{},
		&models.MemoryCost{},
		&models.StorageCost{},
		&models.IngressCost{},
		&models.MonthlyCost{},
	)
	if err != nil {
		return err
	}

	return nil
}
