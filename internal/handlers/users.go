package handlers

import (
	"CryptoPortfolio/internal/models"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var UserAPI models.UserAPI

	if err := json.Unmarshal(b, &UserAPI); err != nil {
		http.Error(w, errors.New("BadRequest").Error(), http.StatusBadRequest)
		return
	}

	token, err := h.app.CreateUser(UserAPI, h.config.Salt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	log.Print(token)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("JWT " + token))
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var UserAPI models.UserAPI

	if err := json.Unmarshal(b, &UserAPI); err != nil {
		http.Error(w, errors.New("BadRequest").Error(), http.StatusBadRequest)
		return
	}

	token, err := h.app.AuthUser(UserAPI, h.config.Salt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", "Bearer "+token)
	log.Print(token)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JWT " + token))
}
