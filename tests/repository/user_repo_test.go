package repository_test

import (
	"context"
	"testing"

	"game-platform/internal/model"
	"game-platform/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// MockUserRepository для тестирования
type MockUserRepository struct {
	users map[string]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*model.User),
	}
}

func (m *MockUserRepository) GetBySID(ctx context.Context, sid string) (*model.User, error) {
	user, exists := m.users[sid]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *MockUserRepository) Create(ctx context.Context, sid, email, name, gender, department, position, photoURL string) (*model.User, error) {
	user := &model.User{
		SID:        sid,
		Email:      email,
		Name:       name,
		Gender:     gender,
		Department: department,
		Position:   position,
		PhotoURL:   photoURL,
	}
	m.users[sid] = user
	return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	if _, exists := m.users[user.SID]; !exists {
		return nil
	}
	m.users[user.SID] = user
	return nil
}

func (m *MockUserRepository) GetUserGroups(ctx context.Context, sid string) ([]string, error) {
	return []string{"games", "tournaments"}, nil
}

func TestMockUserRepository_GetBySID(t *testing.T) {
	mockRepo := NewMockUserRepository()

	// Тест: пользователь не найден
	user, err := mockRepo.GetBySID(context.Background(), "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}

	// Тест: пользователь найден
	_, err = mockRepo.Create(context.Background(), "emp_12345", "test@test.com", "Test User", "male", "IT", "Dev", "")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err = mockRepo.GetBySID(context.Background(), "emp_12345")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Error("Expected user, got nil")
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", user.Name)
	}
}

func TestMockUserRepository_Create(t *testing.T) {
	mockRepo := NewMockUserRepository()

	user, err := mockRepo.Create(
		context.Background(),
		"emp_67890",
		"petrov@test.com",
		"Петров Петр",
		"male",
		"HR",
		"Manager",
		"https://s3.yakbson.digital/avatars/emp_67890.jpg",
	)

	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.SID != "emp_67890" {
		t.Errorf("Expected SID emp_67890, got %s", user.SID)
	}

	if user.Email != "petrov@test.com" {
		t.Errorf("Expected email petrov@test.com, got %s", user.Email)
	}

	if user.Department != "HR" {
		t.Errorf("Expected department HR, got %s", user.Department)
	}
}

func TestMockUserRepository_Update(t *testing.T) {
	mockRepo := NewMockUserRepository()

	// Создаем пользователя
	_, err := mockRepo.Create(
		context.Background(),
		"emp_11111",
		"sidorov@test.com",
		"Сидоров Сидор",
		"male",
		"IT",
		"Senior Dev",
		"",
	)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Обновляем пользователя
	user, err := mockRepo.GetBySID(context.Background(), "emp_11111")
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	user.Position = "Lead Developer"
	user.Department = "Engineering"

	err = mockRepo.Update(context.Background(), user)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	// Проверяем обновление
	updatedUser, err := mockRepo.GetBySID(context.Background(), "emp_11111")
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}

	if updatedUser.Position != "Lead Developer" {
		t.Errorf("Expected position 'Lead Developer', got %s", updatedUser.Position)
	}

	if updatedUser.Department != "Engineering" {
		t.Errorf("Expected department 'Engineering', got %s", updatedUser.Department)
	}
}

func TestMockUserRepository_GetUserGroups(t *testing.T) {
	mockRepo := NewMockUserRepository()

	groups, err := mockRepo.GetUserGroups(context.Background(), "emp_12345")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	hasGames := false
	hasTournaments := false
	for _, group := range groups {
		if group == "games" {
			hasGames = true
		}
		if group == "tournaments" {
			hasTournaments = true
		}
	}

	if !hasGames {
		t.Error("Expected 'games' group")
	}

	if !hasTournaments {
		t.Error("Expected 'tournaments' group")
	}
}

// Тесты для реального репозитория (требуют PostgreSQL)
func TestUserRepository_GetBySID(t *testing.T) {
	// Этот тест требует работающую PostgreSQL базу
	// Пропускаем если база недоступна
	dbURL := "postgresql://game:***@localhost:5432/game_platform?sslmode=disable"
	
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to database: %v", err)
		return
	}
	defer dbPool.Close()

	repo := repository.NewUserRepository(dbPool)

	user, err := repo.GetBySID(context.Background(), "emp_12345")
	if err != nil {
		t.Skipf("Skipping test: user not found: %v", err)
		return
	}

	if user == nil {
		t.Skip("Skipping test: no test user in database")
		return
	}

	if user.SID != "emp_12345" {
		t.Errorf("Expected SID emp_12345, got %s", user.SID)
	}
}
