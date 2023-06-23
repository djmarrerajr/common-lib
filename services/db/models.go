package db

// TODO: refactor out to a non-library package
import "github.com/google/uuid"

type Account struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Balance float32   `gorm:"type:decimal"`
}
