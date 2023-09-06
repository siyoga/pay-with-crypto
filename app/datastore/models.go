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
		Name          string                `json:"name" gorm:"type:string;unique"`
		Image         string                `json:"image" gorm:"type:string"`
		Password      string                `json:"-" gorm:"type:string"`
		Email         string                `json:"email" gorm:"type:string;not null"`
		LinkToCompany string                `json:"linkToCompany" gorm:"type:string"`
		ViaGoogle     bool                  `json:"viaGoogle" gorm:"default:false;not null"`
		CreatedAt     time.Time             `json:"createdAt" gorm:"type:time"`
		UpdateAt      time.Time             `json:"updatedAt" gorm:"type:time"`
		DeletedAt     time.Time             `json:"deletedAt" gorm:"type:time"`
		Cards         []Card                `json:"cards" gorm:"foreignKey:CompanyOwner"`
		IsDel         soft_delete.DeletedAt `json:"is_del" gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
		CreatedTags   []Tag                 `json:"created_tags" gorm:"foreignKey:CreatorID"`
		RefreshToken  RefreshToken          `json:"refresh_token" gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}

	Card struct {
		ID            uuid.UUID             `json:"id" gorm:"type:uuid"`
		CompanyOwner  uuid.UUID             `json:"companyOwner" gorm:"type:uuid"`
		Name          string                `json:"name" gorm:"type:string"`
		LogoLink      string                `json:"logoLink" gorm:"type:string"`
		LinkToWebsite string                `json:"linkToWebsite" gorm:"type:string"`
		Price         string                `json:"price" gorm:"type:string"`
		Description   string                `json:"description" gorm:"type:string"`
		Approved      string                `json:"approved" gorm:"type:string"`
		Tags          pq.StringArray        `json:"tags" gorm:"type:text[]"`
		Views         int                   `json:"views" gorm:"type:int"`
		DeletedAt     time.Time             `json:"deletedAt" gorm:"type:time"`
		IsDel         soft_delete.DeletedAt `json:"isDel" gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
	}

	Admin struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid"`
		Name      string    `json:"name" gorm:"type:string;unique"`
		FirstName string    `json:"firstName" gorm:"type:string"`
		LastName  string    `json:"lastName" gorm:"type:string"`
		Password  string    `json:"password" gorm:"type:string"`
		CreatedAt time.Time //add ``
		UpdatedAt time.Time //add ``
	}

	Tag struct {
		ID        uuid.UUID `json:"id" gorm:"type:uuid"`
		Name      string    `json:"name" gorm:"type:string"`
		CreatorID uuid.UUID `json:"creatorId" gorm:"type:uuid"`
		Approved  string    `json:"approved" gorm:"type:string"`
	}

	RefreshToken struct {
		CompanyID uuid.UUID `json:"companyId" gorm:"type:uuid"`
		Token     string    `json:"token" gorm:"type:string"`
	}
)
