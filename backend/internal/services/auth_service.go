package services

import (
	"errors"
	"finara-backend/internal/config"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"finara-backend/internal/utils"
)

type AuthService interface {
	Register(email, password, fullName string, role models.UserRole) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, config *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		config:   config,
	}
}

func (s *authService) Register(email, password, fullName string, role models.UserRole) (*models.User, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil && existingUser.ID > 0 {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:    email,
		Password: hashedPassword,
		FullName: fullName,
		Role:     role,
		IsActive: true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return "", nil, errors.New("user account is inactive")
	}

	// Verify password
	if !utils.CheckPassword(password, user.Password) {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(
		user.ID,
		user.Email,
		string(user.Role),
		user.CompanyID,
		s.config.JWTSecret,
		s.config.JWTExpiration,
	)

	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}