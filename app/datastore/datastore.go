package datastore

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Datastore *DatastoreT
)

func init() {
	Datastore = new(DatastoreT)
}

func New(parameters string) {
	db, err := gorm.Open(postgres.Open(parameters), &gorm.Config{})

	if err != nil {
		logrus.Error("Error in New.", err)
	}

	db.AutoMigrate(&Company{}, &Card{}, &RefreshToken{}, &Admin{}, &Tag{})

	Datastore = &DatastoreT{db}
}
