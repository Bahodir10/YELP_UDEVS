package service

import (
	"errors"
	"YALP/internal/repository"
	"YALP/internal/util"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"YALP/internal/domain"
)
type UserService interface {
	Register(email, password, name string) (map[string]string, error)
	Login(email, password string) (string, error)
}
type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(r repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: r, jwtSecret: jwtSecret}
}

// Register a new user
func (s *userService) Register(email, password, name string) (map[string]string, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create a user object
	newUser := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	// Save the user and get the generated user ID
	userID, err := s.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	// Return the user details along with the user ID
	return map[string]string{
		"email": email,
		"name":  name,
		"id":    fmt.Sprintf("%d", userID),
	}, nil
}



// Login an existing user
func (s *userService) Login(email, password string) (string, error) {
	// Get the user by email
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := util.GenerateJWT(int(user.ID), s.jwtSecret)  // Cast user.ID to int
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

