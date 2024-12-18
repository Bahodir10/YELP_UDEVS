package domain

import "time"

type User struct {
	ID               int64     `db:"id"`
	Email            string    `db:"email"`
	Password         string    `db:"password"`
	Name             string    `db:"name"`
	ClaimedBusinesses []int64  `db:"claimed_businesses"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
