package main

import (
	"log"
	"os"
	"scrapper2/config"
	"scrapper2/pkg/csvwriter"
	"scrapper2/pkg/scraper"
)

func main() {
	// Load config
	configFile := "config.json"
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		configFile = "config.json.example"
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	// Config CSV
	writer, file, err := csvwriter.NewCSVWriter("bolsas_pub_" + cfg.Year + ".csv")
	if err != nil {
		log.Fatal("Error creating CSV file:", err)
	}
	defer file.Close()
	defer writer.Flush()

	// Create and run scraper
	s := scraper.NewScraper(cfg, writer)
	s.Run()
}
