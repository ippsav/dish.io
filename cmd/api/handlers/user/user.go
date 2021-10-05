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
	FindUser(ctx context.Context, email, username, password string) (*domain.User, error)
}

type UserHandler struct {
	Service Service
	Log     *zerolog.Logger
}
type BodyRequestJSON struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var CookieName = "qid"

func (h *UserHandler) RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ctx := r.Context()
	value := &BodyRequestJSON{}
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
		return
	}
	if value.Password == "" || value.Username == "" || value.Email == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
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

func (h *UserHandler) LoginHandler(rw http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ctx := r.Context()
	value := &BodyRequestJSON{}

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
		return
	}
	if value.Username == "" && value.Email == "" || value.Password == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}
	user, err := h.Service.FindUser(ctx, value.Email, value.Username, value.Password)
	if err != nil {
		h.Log.Err(err).Msg("could not find user")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    user.ID,
		Path:     "/",
		MaxAge:   3 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(user)
	if err != nil {
		h.Log.Err(err).Msgf("could not encode %v into response writer", user)
		return
	}
	h.Log.Info().Msgf("the request took %v to resolve", time.Since(now))
}
