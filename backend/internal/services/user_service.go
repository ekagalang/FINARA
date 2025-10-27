package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
)

type UserService interface {
	GetUserByID(id uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	UpdateUser(id uint, updates map[string]interface{}) (*models.User, error)
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	if fullName, ok := updates["full_name"].(string); ok {
		user.FullName = fullName
	}
	if role, ok := updates["role"].(string); ok {
		user.Role = models.UserRole(role)
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		user.IsActive = isActive
	}
	// TAMBAHKAN INI - untuk update company_id
	if companyID, ok := updates["company_id"]; ok {
		switch v := companyID.(type) {
		case float64:
			user.CompanyID = uint(v)
		case int:
			user.CompanyID = uint(v)
		case uint:
			user.CompanyID = v
		}
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}