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
		Cards        []Card    `json:"card_id" gorm:"foreignKey:UserID"`
	}

	Card struct {
		UserID uuid.UUID `json:"id" gorm:"type:uuid"`
		Name   string    `json:"name" gorm:"type:string"`
	}
)
