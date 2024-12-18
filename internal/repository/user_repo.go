package repository

import (
	"database/sql"
	"errors"
	"YALP/internal/domain"
)

// UserRepository interface defines methods for user operations
type UserRepository interface {
	CreateUser(u *domain.User) (int64, error)
	GetByEmail(email string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	UpdateUser(u *domain.User) error
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository initializes a new UserRepository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUser inserts a new user into the database and returns the generated ID
func (r *userRepository) CreateUser(u *domain.User) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id
	`, u.Email, u.Password, u.Name).Scan(&id)
	if err != nil {
		return 0, errors.New("failed to create user: " + err.Error())
	}
	return id, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`
		SELECT id, email, password, name, claimed_businesses, created_at, updated_at 
		FROM users WHERE email = $1
	`, email).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.ClaimedBusinesses, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user: " + err.Error())
	}
	return &u, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(`
		SELECT id, email, password, name, claimed_businesses, created_at, updated_at 
		FROM users WHERE id = $1
	`, id).Scan(&u.ID, &u.Email, &u.Password, &u.Name, &u.ClaimedBusinesses, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to retrieve user: " + err.Error())
	}
	return &u, nil
}

// UpdateUser updates an existing user's details
func (r *userRepository) UpdateUser(u *domain.User) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET email=$1, password=$2, name=$3, claimed_businesses=$4, updated_at=NOW() 
		WHERE id=$5
	`, u.Email, u.Password, u.Name, u.ClaimedBusinesses, u.ID)
	if err != nil {
		return errors.New("failed to update user: " + err.Error())
	}
	return nil
}
