package postgres

import (
	"context"
	"github.com/dish.io/internal/domain"
	"github.com/pkg/errors"
)

// CreateUser : Creating user from
func (s *Store) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	statement, err := s.DB.PrepareContext(ctx, "INSERT INTO users(username,email,passwordHash) values($1,$2,$3) RETURNING id,created_at")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare insert statement")
	}
	err = statement.QueryRowContext(ctx, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not insert user into db")
	}
	return user, nil
}

func (s *Store) FindUserByID(ctx context.Context, id string) (*domain.User, error) {
	statement, err := s.DB.PrepareContext(ctx, "SELECT * from users where id=$1")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}
	user := &domain.User{}
	err = statement.QueryRowContext(ctx, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not select from users db")
	}
	return user, nil
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	statement, err := s.DB.PrepareContext(ctx, "SELECT * from users where email=$1")
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare select statement")
	}
	user := &domain.User{}
	err = statement.QueryRowContext(ctx, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, errors.Wrap(err, "could not select from users db")
	}
	return user, nil
}

func (s *Store) CheckUserExists(ctx context.Context, username, email string) (bool, error) {
	var count int
	statement, err := s.DB.PrepareContext(ctx, "SELECT COUNT(*) FROM users WHERE username=$1 OR email=$2")
	if err != nil {
		return false, errors.Wrap(err, "could not prepare select statement")
	}
	err = statement.QueryRowContext(ctx, username, email).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "could not select from users table")
	}
	return count != 0, nil
}
