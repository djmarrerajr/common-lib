package db

import "github.com/google/uuid"

type Account struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Balance float32   `gorm:"type:decimal"`
}
