package service

import (
	"YALP/internal/domain"
	"YALP/internal/repository"
	"errors"
)

type ReviewService interface {
	CreateReview(businessID, userID int64, rating int, comment string) (*domain.Review, error)
	ListReviewsForBusiness(businessID int64) ([]domain.Review, error)
}

type reviewService struct {
	reviewRepo   repository.ReviewRepository
	businessRepo repository.BusinessRepository
}

func NewReviewService(r repository.ReviewRepository, b repository.BusinessRepository) ReviewService {
	return &reviewService{reviewRepo: r, businessRepo: b}
}

func (s *reviewService) CreateReview(businessID, userID int64, rating int, comment string) (*domain.Review, error) {
	b, err := s.businessRepo.FindByID(businessID)
	if err != nil || b == nil {
		return nil, errors.New("business does not exist")
	}
	rv := &domain.Review{
		BusinessID: businessID,
		UserID:     userID,
		Rating:     rating,
		Comment:    comment,
	}
	id, err := s.reviewRepo.Create(rv)
	if err != nil {
		return nil, err
	}
	rv.ID = id
	return rv, nil
}

func (s *reviewService) ListReviewsForBusiness(businessID int64) ([]domain.Review, error) {
	return s.reviewRepo.ListByBusinessID(businessID)
}
