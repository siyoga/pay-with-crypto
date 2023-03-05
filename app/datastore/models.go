package datastore

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type (
	DatabaseConfig struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Database string `json:"database"`
	}

	DatastoreT struct {
		*gorm.DB
	}

	All interface {
		User
	}

	User struct {
		ID           uuid.UUID `json:"id" gorm:"type:uuid"`
		Company_Name string    `json:"company_name" gorm:"type:string"`
		Password     string    `json:"password" gorm:"type:string"`
		Cards        []Card    `json:"cards" gorm:"foreignKey:UserID"`
	}

	Card struct {
		UserID      uuid.UUID `json:"id" gorm:"type:uuid"`
		Name        string    `json:"name" gorm:"type:string"`
		LinkToProd  string    `json:"linkToProd" gorm:"type:string"`
		Price       string    `json:"price" gorm:"type:string"`
		Description string    `json:"description" gorm:"type:string"`
		Tags        []string  `json:"tags" gorm:"type:text[]"`
	}

	RefreshToken struct {
		Token string `json:"token" gorm:"string"`
	}
)
