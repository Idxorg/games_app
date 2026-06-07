package repository_test

import (
	"context"
	"testing"

	"game-platform/tests/mocks"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestMockUserRepo_BasicOperations(t *testing.T) {
	mockRepo := mocks.NewMockUserRepo()

	// Test: user not found
	user, err := mockRepo.GetBySID(context.Background(), "nonexistent")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user != nil {
		t.Errorf("Expected nil user, got %v", user)
	}

	// Test: create and retrieve user
	_, err = mockRepo.Create(context.Background(), "emp_12345", "test@test.com", "Test User", "male", "IT", "Dev", "")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user, err = mockRepo.GetBySID(context.Background(), "emp_12345")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("Expected user, got nil")
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got %s", user.Name)
	}

	// Test: update user
	user.Position = "Lead Developer"
	user.Department = "Engineering"
	err = mockRepo.Update(context.Background(), user)
	if err != nil {
		t.Errorf("Failed to update user: %v", err)
	}

	updatedUser, err := mockRepo.GetBySID(context.Background(), "emp_12345")
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}
	if updatedUser.Position != "Lead Developer" {
		t.Errorf("Expected position 'Lead Developer', got %s", updatedUser.Position)
	}
}

func TestMockUserRepo_GetUserGroups(t *testing.T) {
	mockRepo := mocks.NewMockUserRepo()

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

func TestMockUserRepo_GetUserGroups_AlwaysReturns(t *testing.T) {
	mockRepo := mocks.NewMockUserRepo()

	groups, err := mockRepo.GetUserGroups(context.Background(), "unknown")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(groups) == 0 {
		t.Errorf("Mock should return default groups")
	}
}

// Integration test — requires real PostgreSQL
func TestUserRepository_Integration(t *testing.T) {
	dbURL := "postgresql://game:***@localhost:5432/game_platform?sslmode=disable"

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("Skipping test: cannot connect to database: %v", err)
		return
	}
	defer dbPool.Close()

	// This test only runs if a real database is available.
	// The real repository is tested here.
	_ = dbPool // TODO: add actual queries once DB is set up
	t.Skip("Integration test: requires configured PostgreSQL instance")
}
