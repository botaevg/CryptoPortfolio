package servise

import (
	"CryptoPortfolio/internal/ExternalService"
	"CryptoPortfolio/internal/models"
	"CryptoPortfolio/internal/repositories"
	"log"
)

type Portfolio struct {
	storage         repositories.Storage
	externalService ExternalService.ExternalService
}

func NewPortfolio(storage repositories.Storage, externalService ExternalService.ExternalService) Portfolio {
	return Portfolio{
		storage:         storage,
		externalService: externalService,
	}
}

func (p Portfolio) CreatePortfolio(portfolio models.PortfolioAPI, username uint) error {
	var Portfolio models.Portfolio
	Portfolio.Name = portfolio.Name
	Portfolio.UserID = username
	err := p.storage.CreatePortfolio(Portfolio, username)
	return err
}

func (p Portfolio) GetPortfolio(username uint) ([]models.PortfolioAPI, error) {

	Portfolio, err := p.storage.GetPortfolio(username)

	var portfolioAPI []models.PortfolioAPI
	for _, v := range Portfolio {
		x := models.PortfolioAPI{Name: v.Name, ID: v.ID}
		portfolioAPI = append(portfolioAPI, x)
	}
	return portfolioAPI, err
}

func (p Portfolio) AddCoin(coinAPI models.CoinAPI, id uint) error {

	var Coin models.Coin
	Coin.Name = coinAPI.Name
	Coin.PricePurchase = coinAPI.PricePurchase
	Coin.Amount = coinAPI.Amount

	err := p.storage.AddCoin(Coin, id)
	return err
}

func (p Portfolio) OwnerPortfolio(userid uint, id uint) (bool, error) {
	return p.storage.OwnerPortfolio(userid, id)
}

func (p Portfolio) OpenPortfolio(id uint) ([]models.CoinGroup, error) {
	CoinGroup, err := p.storage.OpenPortfolio(id)
	for i, v := range CoinGroup {
		CoinGroup[i].PriceNow, err = p.externalService.GetCurrency(v.Name)
		if err != nil {
			log.Print("getcur error")
			log.Print(err)
		}
		CoinGroup[i].TotalValue = v.Total * CoinGroup[i].PriceNow
	}
	return CoinGroup, err
}
