package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseUUID struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

func (b *BaseUUID) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}
