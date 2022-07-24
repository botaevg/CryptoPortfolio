package repositories

import (
	"CryptoPortfolio/internal/models"
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

type Storage interface {
	Ping(ctx context.Context) error
	CreateUser(user models.User) (uint, error)
	GetUser(user models.User) (uint, error)
	CreatePortfolio(portfolio models.Portfolio, username uint) error
	GetPortfolio(username uint) ([]models.Portfolio, error)
	AddCoin(coin models.Coin, portfolioid uint) error
	OwnerPortfolio(userid uint, id uint) (bool, error)
	OpenPortfolio(id uint) ([]models.CoinGroup, error)
}

type StorageDB struct {
	ConDB *pgx.Conn
}

func NewDB(con *pgx.Conn) *StorageDB {
	return &StorageDB{
		ConDB: con,
	}
}

func (s StorageDB) Ping(ctx context.Context) error {
	if err := s.ConDB.Ping(ctx); err != nil {
		log.Print("ping error")
		return err
	}
	return nil
}

func (s StorageDB) CreateUser(user models.User) (uint, error) {
	q := `INSERT INTO users (username, password)
	VALUES ($1, $2) RETURNING id;`

	log.Print(user.Username, user.Password)

	rows, err := s.ConDB.Query(context.Background(), q, user.Username, user.Password)
	defer rows.Close()
	if err != nil {
		log.Print("Запись не создана")
		return 0, err
	}
	var id uint
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Print("Запись не создана")
			return 0, err
		}
	}

	log.Print("Запись создана")
	log.Print(id)
	return id, nil
}

func (s StorageDB) GetUser(user models.User) (uint, error) {
	q := `SELECT id FROM users WHERE username = $1 and password = $2;`

	row, err := s.ConDB.Query(context.Background(), q, user.Username, user.Password)
	defer row.Close()
	log.Print("query")
	log.Print(user.Username, user.Password)
	if err != nil {
		log.Print("Get User")
		log.Print(err)
		return 0, err
	}
	var id uint
	if row.Next() {
		err := row.Scan(&id)
		if err != nil {
			log.Print(err)
			return 0, err
		}
		log.Print("user found")
		return id, nil
	}
	log.Print("user no found")
	return 0, nil
}

func (s StorageDB) CreatePortfolio(portfolio models.Portfolio, userid uint) error {
	q := `INSERT INTO portfolio (userid, name)
	VALUES ($1, $2);`

	log.Print(portfolio.Name, userid)

	_, err := s.ConDB.Exec(context.Background(), q, userid, portfolio.Name)
	if err != nil {
		log.Print("Запись не создана")
		log.Print(err)
		return err
	}
	log.Print("Запись создана")
	return nil
}

func (s StorageDB) GetPortfolio(username uint) ([]models.Portfolio, error) {
	q := `SELECT id, name FROM portfolio WHERE userid = $1`

	rows, err := s.ConDB.Query(context.Background(), q, username)
	defer rows.Close()
	log.Print("storage rows")

	var Portfolio []models.Portfolio
	for rows.Next() {
		x := models.Portfolio{}
		err := rows.Scan(&x.ID, &x.Name)
		if err != nil {
			return []models.Portfolio{}, err
		}
		Portfolio = append(Portfolio, x)
	}
	if rows.Err() != nil {
		return []models.Portfolio{}, err
	}
	log.Print(Portfolio)
	return Portfolio, err
}

func (s StorageDB) AddCoin(coin models.Coin, portfolioid uint) error {
	q := `INSERT INTO coinaction (portfolioid, name, amount, pricepurchase)
	VALUES ($1, $2, $3, $4);`
	_, err := s.ConDB.Exec(context.Background(), q, portfolioid, coin.Name, coin.Amount, coin.PricePurchase)
	if err != nil {
		log.Print("Запись не создана")
		return err
	}
	log.Print("Запись создана")
	return err
}

func (s StorageDB) OwnerPortfolio(userid uint, id uint) (bool, error) {
	q := `select * from portfolio where userid = $1 and id = $2;`
	rows, err := s.ConDB.Query(context.Background(), q, userid, id)
	defer rows.Close()
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (s StorageDB) OpenPortfolio(id uint) ([]models.CoinGroup, error) {
	q := `select name, sum(amount) as total, sum(amount*pricepurchase)/sum(amount) as averageprice, sum(amount*pricepurchase) as totalinvested  from coinaction
			where portfolioid = $1 group by name;`
	rows, err := s.ConDB.Query(context.Background(), q, id)
	defer rows.Close()
	if err != nil {
		return []models.CoinGroup{}, err
	}

	var CoinGroup []models.CoinGroup

	for rows.Next() {
		x := models.CoinGroup{}
		err := rows.Scan(&x.Name, &x.Total, &x.AveragePrice, &x.TotalInvested)
		if err != nil {
			return []models.CoinGroup{}, err
		}
		CoinGroup = append(CoinGroup, x)
	}

	err = rows.Err()
	if err != nil {
		return []models.CoinGroup{}, err
	}

	return CoinGroup, nil
}
