package db

import (
	"log"
	"strings"

	"github.com/yardnsm/leeches/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	pgSuffix = "pg:"
)

func CreateDatabase(databsePath string) *gorm.DB {
	var dialector gorm.Dialector

	if strings.HasPrefix(databsePath, pgSuffix) {
		dialector = postgres.Open(strings.TrimPrefix(databsePath, pgSuffix))
	} else {
		dialector = sqlite.Open(databsePath)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.ChargeRequest{},
		&model.ChargeMessage{},
	)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
