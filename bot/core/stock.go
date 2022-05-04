package core

import (
	"bot/messaging"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	Client HTTPClient
)

type StockRequest struct {
	Code   string `json:"code"`
	RoomId string `json:"roomId"`
}

type StockResponse struct {
	Symbol string `json:"symbol"`
	Date   string `json:"date"`
	Time   string `json:"time"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

// HTTPClient interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func init() {
	Client = &http.Client{}
}

func GetStockQuote(code string) (string, error) {
	fmt.Println(code)
	if code == "" {
		return "", errors.New("invalid code")
	}
	data, err := ReadCSVFromUrl(messaging.STOCK_URL + code)
	if err != nil {
		return "", errors.New("error parsing CSV from URL")
	}

	dataFieldRows := data[1]
	stooqResponse := &StockResponse{
		Symbol: dataFieldRows[0],
		Close:  dataFieldRows[6],
	}

	var msg string
	if stooqResponse.Close != "N/D" {
		msg = fmt.Sprintf("%s quote is %v per share.", stooqResponse.Symbol, stooqResponse.Close)
	} else {
		msg = "Could not get stock quote."
	}

	return msg, nil
}

func ReadCSVFromUrl(url string) ([][]string, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := Client.Do(request)

	if err != nil {
		log.Print(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid response status code")
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ';'
	reader.LazyQuotes = true
	data, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var rows [][]string
	for _, e := range data {
		rows = append(rows, strings.Split(e[0], ","))
	}

	return rows, nil
}
