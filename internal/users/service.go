package users

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameExists  = errors.New("username already exists")
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidUserID   = errors.New("invalid user id")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	user *User,
) error {
	if user == nil {
		return errors.New("user is required")
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	if user.Username == "" {
		return ErrInvalidUsername
	}

	if user.Email == "" {
		return ErrInvalidEmail
	}

	existingEmail, err := s.repository.FindByEmail(
		ctx,
		user.Email,
	)

	if err == nil && existingEmail != nil {
		return ErrEmailExists
	}

	existingUsername, err := s.repository.FindByUsername(
		ctx,
		user.Username,
	)

	if err == nil && existingUsername != nil {
		return ErrUsernameExists
	}

	return s.repository.Create(ctx, user)
}

func (s *Service) GetByID(
	ctx context.Context,
	id uint64,
) (*User, error) {
	if id == 0 {
		return nil, ErrInvalidUserID
	}

	user, err := s.repository.FindByID(ctx, id)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *Service) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, ErrInvalidEmail
	}

	user, err := s.repository.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *Service) GetByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	username = strings.TrimSpace(username)

	if username == "" {
		return nil, ErrInvalidUsername
	}

	user, err := s.repository.FindByUsername(
		ctx,
		username,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *Service) Update(
	ctx context.Context,
	user *User,
) error {
	if user == nil || user.ID == 0 {
		return ErrInvalidUserID
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	if user.Username == "" {
		return ErrInvalidUsername
	}

	if user.Email == "" {
		return ErrInvalidEmail
	}

	return s.repository.Update(ctx, user)
}

func (s *Service) Delete(
	ctx context.Context,
	id uint64,
) error {
	if id == 0 {
		return ErrInvalidUserID
	}

	if _, err := s.repository.FindByID(ctx, id); err != nil {
		return ErrUserNotFound
	}

	return s.repository.Delete(ctx, id)
}
