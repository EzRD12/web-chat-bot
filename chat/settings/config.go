package settings

import (
	"encoding/json"
	"log"
	"os"
)

var Configuration Config

type Config struct {
	RabbitMQ struct {
		User           string `json:"user"`
		Pass           string `json:"pass"`
		Host           string `json:"host"`
		Port           string `json:"port"`
		GetStockQueue  string `json:"getStockQueue"`
		SendStockQueue string `json:"sendStockQueue"`
	}

	AppConfig struct {
		SecretKey string `json:"secretKey"`
	}
}

func GetConfig() *Config {
	var cfg Config
	readFile(&cfg)

	Configuration = cfg

	return &cfg
}

func readFile(cfg *Config) {
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Fatal(err)
	}
}
