package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BaseURL      string `json:"baseURL"`
	LoginURL     string `json:"loginURL"`
	AuthURL      string `json:"authURL"`
	ListPageURL  string `json:"listPageURL"`
	DwrApiURL    string `json:"dwrApiURL"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Year         string `json:"year"`
	UserAgent    string `json:"userAgent"`
	Timeout      int    `json:"timeout"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
