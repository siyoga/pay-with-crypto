package datastore

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
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
		Company | RefreshToken | Card | Admin | Tag
	}

	Authable interface {
		Company | Admin
	}

	Company struct {
		ID            uuid.UUID             `json:"id" gorm:"type:uuid"`
		Username      string                `json:"username" gorm:"type:string;unique"`
		Image         string                `json:"image" gorm:"type:string"`
		Password      string                `json:"-" gorm:"type:string"`
		Mail          string                `json:"mail" gorm:"type:string"`
		LinkToCompany string                `json:"linkToCompany" gorm:"type:string"`
		ViaGoogle     bool                  `json:"viaGoogle" gorm:"default:false"`
		CreatedAt     time.Time             `json:"createdAt" gorm:"type:time"`
		UpdateAt      time.Time             `json:"updatedAt" gorm:"type:time"`
		DeletedAt     time.Time             `json:"deletedAt" gorm:"type:time"`
		Cards         []Card                `json:"cards" gorm:"foreignKey:CompanyID"`
		IsDel         soft_delete.DeletedAt `json:"is_del" gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
		CreatedTags   []Tag                 `json:"created_tags" gorm:"foreignKey:CreatorID"`
		RefreshToken  RefreshToken          `json:"refresh_token" gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	Card struct {
		ID          uuid.UUID             `json:"id" gorm:"type:uuid"`
		CompanyID   uuid.UUID             `json:"company_id" gorm:"type:uuid"`
		Name        string                `json:"name" gorm:"type:string"`
		Image       string                `json:"image" gorm:"type:string"`
		LinkToProd  string                `json:"linkToProd" gorm:"type:string"`
		Price       string                `json:"price" gorm:"type:string"`
		Description string                `json:"description" gorm:"type:string"`
		Approved    string                `json:"approved" gorm:"type:string"`
		Tags        pq.StringArray        `json:"tags" gorm:"type:text[]"`
		Views       int                   `json:"views" gorm:"type:int"`
		DeletedAt   time.Time             `json:"deletedAt" gorm:"type:time"`
		IsDel       soft_delete.DeletedAt `json:"is_del" gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	}

	Admin struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid"`
		Name      string    `json:"name" gorm:"type:string;unique"`
		FirstName string    `json:"first_name" gorm:"type:string"`
		LastName  string    `json:"last_name" gorm:"type:string"`
		Password  string    `json:"password" gorm:"type:string"`
		CreatedAt time.Time //add ``
		UpdatedAt time.Time //add ``
	}

	Tag struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid"`
		Name      string    `json:"name" gorm:"type:string"`
		CreatorID uuid.UUID `json:"creator_id" gorm:"type:uuid"`
		Approved  string    `json:"approved" gorm:"type:string"`
	}

	RefreshToken struct {
		CompanyID uuid.UUID `json:"company_id" gorm:"type:uuid"`
		Token     string    `json:"token" gorm:"type:string"`
	}
)
