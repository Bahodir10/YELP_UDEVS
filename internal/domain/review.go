package domain

import "time"

type Review struct {
	ID         int64     `db:"id"`
	BusinessID int64     `db:"business_id"`
	UserID     int64     `db:"user_id"`
	Rating     int       `db:"rating"`
	Comment    string    `db:"comment"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
