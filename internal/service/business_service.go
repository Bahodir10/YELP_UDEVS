package service

import (
	"errors"
	"YALP/internal/domain"
	"YALP/internal/repository"
)

type BusinessService interface {
	CreateBusiness(name, category, desc, address, contact string, ownerID int64) (*domain.Business, error)
	ListAll() ([]domain.Business, error)
	Search(keyword string) ([]domain.Business, error)
	GetBusiness(id int64) (*domain.Business, error)
	ClaimBusiness(businessID, userID int64) (*domain.Business, error)
}

type businessService struct {
	businessRepo repository.BusinessRepository
	userRepo     repository.UserRepository
}

func NewBusinessService(b repository.BusinessRepository, u repository.UserRepository) BusinessService {
	return &businessService{businessRepo: b, userRepo: u}
}

func (s *businessService) CreateBusiness(name, category, desc, address, contact string, ownerID int64) (*domain.Business, error) {
	b := &domain.Business{
		Name:        name,
		Category:    category,
		Description: desc,
		Address:     address,
		ContactInfo: contact,
		OwnerID:     &ownerID,
	}
	id, err := s.businessRepo.Create(b)
	if err != nil {
		return nil, err
	}
	b.ID = id
	return b, nil
}

func (s *businessService) ListAll() ([]domain.Business, error) {
	return s.businessRepo.ListAll()
}

func (s *businessService) Search(keyword string) ([]domain.Business, error) {
	return s.businessRepo.Search(keyword)
}

func (s *businessService) GetBusiness(id int64) (*domain.Business, error) {
	return s.businessRepo.FindByID(id)
}

func (s *businessService) ClaimBusiness(businessID, userID int64) (*domain.Business, error) {
	b, err := s.businessRepo.FindByID(businessID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, errors.New("business not found")
	}
	// Check if user is allowed to claim
	u, err := s.userRepo.GetByID(userID)
	if err != nil || u == nil {
		return nil, errors.New("user not found")
	}
	if b.OwnerID != nil {
		return nil, errors.New("business already has an owner")
	}
	b.OwnerID = &userID
	if err := s.businessRepo.Update(b); err != nil {
		return nil, err
	}
	u.ClaimedBusinesses = append(u.ClaimedBusinesses, b.ID)
	if err := s.userRepo.UpdateUser(u); err != nil {
		return nil, err
	}
	return b, nil
}
