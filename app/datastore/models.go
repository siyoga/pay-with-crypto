package datastore

import "gorm.io/gorm"

type (
	DatastoreT struct {
		*gorm.DB
	}
)