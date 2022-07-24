package servise

import (
	"CryptoPortfolio/internal/models"
	"CryptoPortfolio/internal/repositories"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"time"
)

type Auth struct {
	storage   repositories.Storage
	secretkey string
}

type Claims struct {
	jwt.StandardClaims
	Username uint `json:"userid"`
}

func GenerateHashForPass(password string, salt string) string {
	h := hmac.New(sha256.New, ([]byte(salt)))
	log.Print("1")
	log.Print(h)
	h.Write([]byte(password))
	log.Print("2")
	log.Print(h)
	sum := h.Sum(nil)
	log.Print("3")
	log.Print(sum)
	return hex.EncodeToString(sum)
}

func (a Auth) CreateUser(userAPI models.UserAPI, salt string) (string, error) {

	var User models.User
	User.Username = userAPI.Username
	User.Password = GenerateHashForPass(userAPI.Password, salt)

	ID, err := a.storage.CreateUser(User)
	if err != nil {
		return "", err
	}
	User.ID = ID
	log.Print("before token")
	log.Print(User.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: User.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.At(time.Now()),
		},
	})
	tokenSigned, err := token.SignedString([]byte(a.secretkey))
	if err != nil {
		log.Print("tokenSigned")
	}
	return tokenSigned, err
}

func (a Auth) AuthUser(userAPI models.UserAPI, salt string) (string, error) {
	var User models.User
	User.Username = userAPI.Username
	User.Password = GenerateHashForPass(userAPI.Password, salt)
	ID, err := a.storage.GetUser(User)
	User.ID = ID
	log.Print("before token")
	log.Print(User.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username: User.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(60 * time.Minute)),
			IssuedAt:  jwt.At(time.Now()),
		},
	})
	tokenSigned, err := token.SignedString([]byte(a.secretkey))
	if err != nil {
		log.Print("tokenSigned")
	}
	return tokenSigned, err
}

func NewAuth(storage repositories.Storage, secretkey string) Auth {
	return Auth{
		storage:   storage,
		secretkey: secretkey,
	}
}
