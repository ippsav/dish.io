package domain

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

func HashPassword(passwordToHash string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordToHash), 8)
	if err != nil {
		return "", errors.Wrap(err, "could not hash password")
	}
	return string(bytes), nil
}
