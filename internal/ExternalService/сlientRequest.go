package ExternalService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type ExternalService struct {
	ClientKey string
}

type USD struct {
	USD float64
}

func (e ExternalService) GetCurrency(coinName string) (float64, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("https://min-api.cryptocompare.com/data/price?fsym=%s&tsyms=USD", coinName)
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Set("Authorization", "Apikey  "+"b50c8970d698bf376f3812aaf8d82a371865102c6c6ff95f2c8e5ed4373790c8")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(respBody))
	var Currency USD
	err = json.Unmarshal(respBody, &Currency)
	if err != nil {
		log.Print(err)
	}
	log.Print(Currency.USD)
	return Currency.USD, nil
}

func NewES(key string) ExternalService {
	return ExternalService{
		ClientKey: key,
	}
}
