package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseUUID
	Name      string     `gorm:"type:varchar(255);not null" json:"name"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Role      string     `gorm:"type:varchar(50);default:'user'" json:"role"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
