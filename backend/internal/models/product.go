package models

import (
	"time"
)

type Product struct {
	BaseUUID
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Description *string    `gorm:"type:text" json:"description"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock       int        `gorm:"not null;default:0" json:"stock"`
	Category    *string    `gorm:"type:varchar(100)" json:"category"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
