package portfolio

import (
	"CryptoPortfolio/internal/ExternalService"
	"CryptoPortfolio/internal/config"
	"CryptoPortfolio/internal/handlers"
	"CryptoPortfolio/internal/middleapp"
	"CryptoPortfolio/internal/repositories"
	"CryptoPortfolio/internal/servise"
	"CryptoPortfolio/pkg/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func StartApp() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Print("config error")
		return
	}
	postgreCon, err := db.NewClient(cfg.DataBaseDSN)
	if err != nil {
		log.Print("DB connect error")
		return
	}
	var storage repositories.Storage
	storage = repositories.NewDB(postgreCon)

	log.Print(cfg.BaseURL)
	log.Print(cfg.ServerAddress)
	log.Print(cfg.DataBaseDSN)
	authService := servise.NewAuth(storage, cfg.SecretKey)

	externalService := ExternalService.NewES(cfg.CoinKey)
	portfolioService := servise.NewPortfolio(storage, externalService)
	r := chi.NewRouter()

	h := handlers.New(cfg, authService, portfolioService)
	auth := middleapp.NewAuth(cfg.SecretKey)
	r.Use(auth.AuthHeader)
	r.Use(middleware.Logger)

	r.Post("/api/login", h.Login)
	r.Post("/api/register", h.RegisterUser)

	r.Get("/api/portfolio", h.GetPortfolio)
	r.Post("/api/portfolio", h.CreatePortfolio)

	r.Get("/api/portfolio/{id}", h.OpenPortfolio)
	r.Post("/api/portfolio/{id}", h.AddCoin)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
