package models

type User struct {
	ID       uint
	Username string
	Password string
}

type UserAPI struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Portfolio struct {
	ID     uint
	UserID uint
	Name   string
	//TotalInvested float64
}

type PortfolioAPI struct {
	Name string `json:"portfolioname"`
	ID   uint   `json:"id"`
}

type Coin struct {
	ActionID      uint
	PortfolioID   uint
	Name          string
	Amount        float64
	PricePurchase float64
	//TotalInvested int

	//TotalValue    int
	//TotalGain     int
}

type CoinAPI struct {
	Name   string  `json:"coin"`
	Amount float64 `json:"amount"`
	//TotalInvested int
	PricePurchase float64 `json:"pricepurchase"`
	//TotalValue    int
	//TotalGain     int
}

type CoinGroup struct {
	Name          string  `json:"coin"`
	Total         float64 `json:"totalcount"`
	AveragePrice  float64 `json:"averageprice"`
	TotalInvested float64 `json:"totalinvested"`
	PriceNow      float64 `json:"pricenow"`
	TotalValue    float64 `json:"totalvalue"`
}
