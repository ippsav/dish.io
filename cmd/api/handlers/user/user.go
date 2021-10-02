package handlers

import (
	"context"
	"encoding/json"
	"github.com/dish.io/internal/domain"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type Service interface {
	CreateUser(ctx context.Context, email, username, password string) (*domain.User, error)
}

type UserHandler struct {
	Service Service
	Log     *zerolog.Logger
}
type CreateUserRequestJSON struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUserHandler(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ctx := r.Context()
	value := &CreateUserRequestJSON{}
	switch r.Header.Get("content-type") {
	case "application/json":
		err := json.NewDecoder(r.Body).Decode(&value)
		if err != nil {
			h.Log.Err(err).Msg("Could not parse body")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		rw.WriteHeader(http.StatusNotAcceptable)
	}
	user, err := h.Service.CreateUser(ctx, value.Email, value.Username, value.Password)
	if err != nil {
		h.Log.Err(err).Msg("Could not create user")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		h.Log.Err(err).Msgf("could not encode %v into response writer", user)
		return
	}
	h.Log.Info().Msgf("the request took %v to resolve", time.Since(now))
}
