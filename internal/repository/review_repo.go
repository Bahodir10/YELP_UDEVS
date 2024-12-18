package repository

import (
	"database/sql"
	"YALP/internal/domain"
)

type ReviewRepository interface {
	Create(r *domain.Review) (int64, error)
	ListByBusinessID(businessID int64) ([]domain.Review, error)
}

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(rv *domain.Review) (int64, error) {
	var id int64
	err := r.db.QueryRow(`
		INSERT INTO reviews (business_id, user_id, rating, comment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id
	`, rv.BusinessID, rv.UserID, rv.Rating, rv.Comment).Scan(&id)
	return id, err
}

func (r *reviewRepository) ListByBusinessID(businessID int64) ([]domain.Review, error) {
	rows, err := r.db.Query(`
		SELECT id, business_id, user_id, rating, comment, created_at, updated_at
		FROM reviews
		WHERE business_id=$1
	`, businessID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rv domain.Review
		if err := rows.Scan(&rv.ID, &rv.BusinessID, &rv.UserID, &rv.Rating, &rv.Comment, &rv.CreatedAt, &rv.UpdatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, rv)
	}
	return reviews, nil
}
