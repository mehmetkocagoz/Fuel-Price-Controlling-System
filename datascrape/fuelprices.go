package datascrape

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Get fuel prices from tppd.com.tr
// This function will just return the fuel prices as a document.
func GetFuelPrices() *goquery.Document {
	url := "https://www.tppd.com.tr/en/former-oil-prices?id=35&county=429&StartDate=22.11.2018&EndDate=22.08.2023"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Get request has failed: ", err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("html.Parse has failed: ", err)
	}
	return doc
}

func ScrapeDateAndFuelPrices(doc goquery.Document) {
	doc.Find("table").Each(func(i int, s *goquery.Selection) {
		s.Find("tr").Each(func(i int, s *goquery.Selection) {
			s.Find("td")

		})
	})
}
