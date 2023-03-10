package datastore

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
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
		User | RefreshToken | Card | Admin | Tag
	}

	User struct {
		ID            uuid.UUID `json:"id" gorm:"type:uuid"`
		Company_Name  string    `json:"company_name" gorm:"type:string"`
		Password      string    `json:"password" gorm:"type:string"`
		Mail          string    `json:"mail" gorm:"type:string"`
		LinkToCompany string    `json:"linkToCompany" gorm:"type:string"`
		Cards         []Card    `json:"cards" gorm:"foreignKey:UserID"`
	}

	Card struct {
		ID          uuid.UUID      `json:"id" gorm:"type:uuid"`
		UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid"`
		Name        string         `json:"name" gorm:"type:string"`
		Image       string         `json:"image" gorm:"type:string"`
		LinkToProd  string         `json:"linkToProd" gorm:"type:string"`
		Price       string         `json:"price" gorm:"type:string"`
		Description string         `json:"description" gorm:"type:string"`
		Approved    bool           `json:"approved" gorm:"type:bool"`
		Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	}

	Admin struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid"`
		UserName  string    `json:"username" gorm:"type:string"`
		FirstName string    `json:"first_name" gorm:"type:string"`
		LastName  string    `json:"last_name" gorm:"type:string"`
		Password  string    `json:"password" gorm:"type:string"`
		CreatedAt time.Time //add ``
		UpdatedAt time.Time //add ``
	}

	Tag struct {
		ID      uuid.UUID `json:"id" gorm:"type:uuid"`
		Name    string    `json:"tagname" gorm:"type:string"`
		AdminID uuid.UUID //add ''
	}

	RefreshToken struct {
		Token string `json:"token" gorm:"string"`
	}
)
