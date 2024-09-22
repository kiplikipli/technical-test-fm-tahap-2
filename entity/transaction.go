package entity

import (
	"time"

	guuid "github.com/google/uuid"
)

type Transaction struct {
	ID                  guuid.UUID  `gorm:"primaryKey" json:"id"`
	UserID              guuid.UUID  `json:"user_id"`
	Type                string      `json:"type"`
	Category            string      `json:"category"`
	Amount              int64       `json:"amount"`
	Remarks             string      `json:"remarks"`
	Status              string      `json:"status"`
	BalanceBefore       int64       `json:"balance_before"`
	BalanceAfter        int64       `json:"balance_after"`
	CorrespondingUserID *guuid.UUID `json:"corresponding_user_id"`
	CreatedAt           time.Time   `gorm:"autoCreateTime" json:"created_at" `
	UpdatedAt           time.Time   `gorm:"autoUpdateTime:milli" json:"-"`

	User              User
	CorrespondingUser User
}
