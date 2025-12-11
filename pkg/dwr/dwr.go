package dwr

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"scrapper2/config"
	"scrapper2/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// This function simulates what the DWR JavaScript would do in the browser
func PrepararSessaoDWR(c *colly.Collector, cfg *config.Config) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	fakeDwrSessionId := string(b)

	cookieDwr := &http.Cookie{
		Name:   "DWRSESSIONID",
		Value:  fakeDwrSessionId,
		Path:   "/jupiterweb",
		Domain: "uspdigital.usp.br",
	}
	c.SetCookies(cfg.BaseURL, []*http.Cookie{cookieDwr})

	randToken := strconv.Itoa(rand.Intn(999999))
	return fakeDwrSessionId + "/" + randToken
}

func DispararDWR(c *colly.Collector, writer *csv.Writer, scriptSessionId string, cfg *config.Config) {
	dwrCollector := c.Clone()

	// Apply the same timeout to the clone
	dwrCollector.SetRequestTimeout(time.Duration(cfg.Timeout) * time.Second)

	dwrCollector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "text/plain")
		r.Headers.Set("Origin", "https://uspdigital.usp.br")
		r.Headers.Set("Referer", cfg.ListPageURL)
	})

	dwrCollector.OnResponse(func(r *colly.Response) {
		body := string(r.Body)
		if strings.Contains(body, "error") || strings.Contains(body, "Exception") {
			fmt.Println("❌ Error returned by server:", body[:utils.Min(len(body), 200)])
			return
		}

		ParseDWRResponse(body, writer)
	})

	dwrCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("❌ FATAL error in DWR:", err)
	})

	payload := strings.Join([]string{
		"callCount=1",
		"nextReverseAjaxIndex=0",
		"c0-scriptName=BeneficioBolsaUnificadaControleDWR",
		"c0-methodName=listarBeneficioBolsaUnificada",
		"c0-id=0",
		"c0-param0=string:" + cfg.Year,
		"c0-param1=string:",
		"c0-param2=string:",
		"c0-param3=string:",
		"c0-param4=boolean:false",
		"batchId=1",
		"instanceId=0",
		"page=%2Fjupiterweb%2FbeneficioBolsaUnificadaListar%3Fcodmnu%3D6684",
		"scriptSessionId=" + scriptSessionId,
		"",
	}, "\n")

	err := dwrCollector.PostRaw(cfg.DwrApiURL, []byte(payload))
	if err != nil {
		fmt.Println("❌ Error sending request:", err)
	}
}

func ParseDWRResponse(body string, writer *csv.Writer) {
	start := strings.Index(body, "[{")
	end := strings.LastIndex(body, "}]")

	if start == -1 || end == -1 {
		fmt.Println("⚠️ No records found or empty response.")
		return
	}

	rawArray := body[start : end+2]
	reObj := regexp.MustCompile(`\{[\s\S]*?\}`)
	objetos := reObj.FindAllString(rawArray, -1)

	reTitulo := regexp.MustCompile(`titprjbnf:\s*"(.*?)"`)
	reUnidade := regexp.MustCompile(`nomabvclg:\s*"(.*?)"`)
	reVertente := regexp.MustCompile(`stavteprj:\s*"(.*?)"`)
	reAno := regexp.MustCompile(`anoofebnf:\s*"(.*?)"`)
	reBolsas := regexp.MustCompile(`numbolapr:\s*(\d+)`)

	count := 0
	for _, obj := range objetos {
		ano := utils.UnquoteUnicode(utils.Extract(reAno, obj))
		unidade := utils.UnquoteUnicode(utils.Extract(reUnidade, obj))
		titulo := utils.UnquoteUnicode(utils.Extract(reTitulo, obj))
		vertente := utils.UnquoteUnicode(utils.Extract(reVertente, obj))
		bolsas := utils.Extract(reBolsas, obj)

		writer.Write([]string{ano, unidade, titulo, vertente, bolsas})
		count++
	}
	writer.Flush()
	fmt.Printf("✅ SUCCESS! %d scholarships saved.\n", count)
}
