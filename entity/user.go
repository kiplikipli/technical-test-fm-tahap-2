package entity

import (
	"time"

	guuid "github.com/google/uuid"
)

type User struct {
	ID          guuid.UUID `gorm:"primaryKey" json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Address     string     `json:"address"`
	PhoneNumber string     `json:"phone_number" gorm:"uniqueIndex"`
	Pin         string     `json:"-"`
	Balance     int64      `json:"balance" gorm:"default:0"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at" `
	UpdatedAt   time.Time  `gorm:"autoUpdateTime:milli" json:"-"`
}
