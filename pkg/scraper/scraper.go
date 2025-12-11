package scraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"scrapper2/config"
	"scrapper2/pkg/dwr"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	cfg       *config.Config
	writer    *csv.Writer
	collector *colly.Collector
}

func NewScraper(cfg *config.Config, writer *csv.Writer) *Scraper {
	c := colly.NewCollector(
		colly.UserAgent(cfg.UserAgent),
		colly.AllowURLRevisit(),
	)

	c.SetRequestTimeout(time.Duration(cfg.Timeout) * time.Second)

	c.WithTransport(&http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(cfg.Timeout) * time.Second,
			KeepAlive: time.Duration(cfg.Timeout) * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	return &Scraper{
		cfg:       cfg,
		writer:    writer,
		collector: c,
	}
}

func (s *Scraper) Run() {
	// Session variables for DWR
	var scriptSessionId string

	// Error handling
	s.collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("❌ Error in request:", err)
		if r.StatusCode != 200 {
			fmt.Println("Status Code:", r.StatusCode)
		}
	})

	// Login
	s.collector.OnHTML("form[action='autenticar']", func(e *colly.HTMLElement) {
		fmt.Println("🔐 Trying to log in...")
		// Small pause to avoid looking too much like a robot
		time.Sleep(1 * time.Second)

		err := s.collector.Post(s.cfg.AuthURL, map[string]string{
			"codpes": s.cfg.User,
			"senusu": s.cfg.Password,
			"Submit": "Entrar",
		})
		if err != nil {
			log.Fatal("Erro no POST de login:", err)
		}
	})

	// Monitor Login and Redirections
	s.collector.OnResponse(func(r *colly.Response) {
		currentURL := r.Request.URL.String()

		if strings.Contains(currentURL, "autenticar") {
			fmt.Println("✅ Login accepted. Initializing session...")
			s.collector.Visit(s.cfg.ListPageURL)

		} else if strings.Contains(currentURL, "beneficioBolsaUnificadaListar") {
			fmt.Println("🛠️ Generating DWR security tokens...")
			scriptSessionId = dwr.PrepararSessaoDWR(s.collector, s.cfg)

			fmt.Println("📡 Sending data request (This may take up to 3 minutes)...")
			dwr.DispararDWR(s.collector, s.writer, scriptSessionId, s.cfg)
		}
	})

	fmt.Println("🚀 Starting crawler (Timeout set to 3 min)...")
	s.collector.Visit(s.cfg.LoginURL)
}
