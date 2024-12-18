package domain

import "time"

type Business struct {
	ID         int64     `db:"id"`
	Name       string    `db:"name"`
	Category   string    `db:"category"`
	Description string   `db:"description"`
	Address    string    `db:"address"`
	ContactInfo string   `db:"contact_info"`
	OwnerID    *int64    `db:"owner_id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
