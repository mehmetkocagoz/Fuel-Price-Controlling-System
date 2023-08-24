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
	var date string
	var diesel string
	doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			// As I know how data table structered I can get the data I want.
			// I'm going to get the date and DIESEL prices only.
			// This function can be improved.
			if i == 0 {
				date = s.Text()
			} else if i == 4 {
				diesel = s.Text()
			}
		})
		fmt.Println(date, diesel)
	})
}
