package user

import (
	"context"
	"github.com/dish.io/internal/domain"
	"github.com/pkg/errors"
)

type Store interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	FindUserByID(ctx context.Context, id string) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CheckUserExists(ctx context.Context, username, email string) (bool, error)
}

type Service struct {
	Store Store
}

func (s *Service) CreateUser(ctx context.Context, email, username, password string) (*domain.User, error) {
	userExists, err := s.Store.CheckUserExists(ctx, username, email)
	if err != nil {
		return nil, err
	}
	if userExists {
		return nil, errors.New("user already has a record in database")
	}
	hashedPassword, err := domain.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	user, err = s.Store.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) FindUser(ctx context.Context, email, username, password string) (*domain.User, error) {
	user := &domain.User{}
	exists, err := s.Store.CheckUserExists(ctx, username, email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	if email != "" {
		user, err = s.Store.FindUserByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
	}
	if username != "" {
		user, err = s.Store.FindUserByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
	}
	passwordCheck := domain.ComparePassword(user.PasswordHash, password)
	if !passwordCheck {
		return nil, nil
	}
	return user, nil
}
