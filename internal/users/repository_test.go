package users

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	_ = godotenv.Load("../../.env")

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		t.Fatal("DATABASE_URL is not configured")
	}

	db, err := gorm.Open(
		postgres.Open(databaseURL),
		&gorm.Config{},
	)

	if err != nil {
		t.Fatalf(
			"failed to connect test database: %v",
			err,
		)
	}

	return db
}

func TestUserRepository(t *testing.T) {
	db := setupTestDB(t)

	repo := NewRepository(db)

	ctx := context.Background()

	user := &User{
		Username:     "avigo_test_user",
		Email:        "avigo_test@example.com",
		PasswordHash: "test_hash",
		Role:         "user",
		IsActive:     true,
	}

	// CREATE
	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf(
			"create user failed: %v",
			err,
		)
	}

	if user.ID == 0 {
		t.Fatal("expected user ID to be generated")
	}

	// Cleanup
	defer func() {
		_ = repo.Delete(ctx, user.ID)
	}()

	// FIND BY ID
	foundByID, err := repo.FindByID(
		ctx,
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"find by ID failed: %v",
			err,
		)
	}

	if foundByID.Email != user.Email {
		t.Fatalf(
			"expected email %s, got %s",
			user.Email,
			foundByID.Email,
		)
	}

	// FIND BY EMAIL
	foundByEmail, err := repo.FindByEmail(
		ctx,
		user.Email,
	)

	if err != nil {
		t.Fatalf(
			"find by email failed: %v",
			err,
		)
	}

	if foundByEmail.Username != user.Username {
		t.Fatalf(
			"expected username %s, got %s",
			user.Username,
			foundByEmail.Username,
		)
	}

	// FIND BY USERNAME
	foundByUsername, err := repo.FindByUsername(
		ctx,
		user.Username,
	)

	if err != nil {
		t.Fatalf(
			"find by username failed: %v",
			err,
		)
	}

	if foundByUsername.ID != user.ID {
		t.Fatalf(
			"expected ID %d, got %d",
			user.ID,
			foundByUsername.ID,
		)
	}

	// UPDATE
	user.Role = "admin"
	user.UpdatedAt = time.Now()

	if err := repo.Update(ctx, user); err != nil {
		t.Fatalf(
			"update user failed: %v",
			err,
		)
	}

	updatedUser, err := repo.FindByID(
		ctx,
		user.ID,
	)

	if err != nil {
		t.Fatalf(
			"find updated user failed: %v",
			err,
		)
	}

	if updatedUser.Role != "admin" {
		t.Fatalf(
			"expected role admin, got %s",
			updatedUser.Role,
		)
	}
}
