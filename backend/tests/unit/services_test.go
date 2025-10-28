package unit

import (
	"finara-backend/internal/models"
	"finara-backend/internal/services"
	"testing"
)

// Mock Repository for testing
type MockUserRepository struct{}

func (m *MockUserRepository) Create(user *models.User) error {
	user.ID = 1
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	if email == "test@example.com" {
		return &models.User{
			BaseModel: models.BaseModel{ID: 1},
			Email:     email,
			FullName:  "Test User",
			Role:      models.RoleAdmin,
		}, nil
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
	return &models.User{
		BaseModel: models.BaseModel{ID: id},
		Email:     "test@example.com",
		FullName:  "Test User",
	}, nil
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	return []models.User{
		{
			BaseModel: models.BaseModel{ID: 1},
			Email:     "test1@example.com",
			FullName:  "User 1",
		},
		{
			BaseModel: models.BaseModel{ID: 2},
			Email:     "test2@example.com",
			FullName:  "User 2",
		},
	}, nil
}

func (m *MockUserRepository) Update(user *models.User) error {
	return nil
}

func (m *MockUserRepository) Delete(id uint) error {
	return nil
}

// Test User Service
func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := services.NewUserService(mockRepo)

	user, err := userService.GetUserByID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID 1, got %d", user.ID)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	if user.FullName != "Test User" {
		t.Errorf("Expected full name Test User, got %s", user.FullName)
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := services.NewUserService(mockRepo)

	users, err := userService.GetAllUsers()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	// Test first user
	if users[0].Email != "test1@example.com" {
		t.Errorf("Expected first user email test1@example.com, got %s", users[0].Email)
	}

	// Test second user
	if users[1].FullName != "User 2" {
		t.Errorf("Expected second user name User 2, got %s", users[1].FullName)
	}
}
