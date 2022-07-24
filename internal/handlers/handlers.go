package handlers

import (
	"CryptoPortfolio/internal/config"
	"CryptoPortfolio/internal/middleapp"
	"CryptoPortfolio/internal/models"
	"CryptoPortfolio/internal/servise"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	config    config.Config
	app       servise.Auth
	portfolio servise.Portfolio
}

func New(cfg config.Config, app servise.Auth, portfolio servise.Portfolio) *handler {
	return &handler{
		config:    cfg,
		app:       app,
		portfolio: portfolio,
	}
}

func (h *handler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(middleapp.Username("username")).(uint)
	log.Print(username)

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var Portfolio models.PortfolioAPI

	if err := json.Unmarshal(b, &Portfolio); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.portfolio.CreatePortfolio(Portfolio, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("portfolio created"))

}

func (h *handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(middleapp.Username("username")).(uint)
	log.Print(username)
	var Portfolio []models.PortfolioAPI

	Portfolio, err := h.portfolio.GetPortfolio(username)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(Portfolio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if string(b) == "null" {
		w.Write([]byte("no portfolios"))
	} else {
		w.Write(b)
	}
}

func (h *handler) OpenPortfolio(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(middleapp.Username("username")).(uint)
	log.Print(username)

	idstr := chi.URLParam(r, "id")
	idint, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := uint(idint)
	owner, err := h.portfolio.OwnerPortfolio(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !owner {
		log.Print("owner false")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var CoinGroup []models.CoinGroup

	CoinGroup, err = h.portfolio.OpenPortfolio(id)

	b, err := json.Marshal(CoinGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
	//id := chi.URLParam(r, "id")

}

func (h *handler) AddCoin(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(middleapp.Username("username")).(uint)
	log.Print(username)

	idstr := chi.URLParam(r, "id")
	idint, err := strconv.Atoi(idstr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := uint(idint)
	owner, err := h.portfolio.OwnerPortfolio(username, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !owner {
		log.Print("owner false")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var coinAPI models.CoinAPI

	err = json.Unmarshal(b, &coinAPI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.portfolio.AddCoin(coinAPI, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("add coin action"))
}
