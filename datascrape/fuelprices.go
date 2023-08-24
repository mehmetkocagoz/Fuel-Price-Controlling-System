package datascrape

import (
	"fmt"
	"mehmetkocagz/database"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FuelPrice struct {
	Date   string
	Diesel string
}

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

func ScrapeDateAndFuelPrices(doc goquery.Document) []FuelPrice {
	var fuelPrices []FuelPrice
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
				// Assume there is no error.
				// We can handle it later.
				diesel = s.Text()
			}
		})
		fuelPrices = append(fuelPrices, FuelPrice{date, diesel})
	})
	return fuelPrices
}

func convertTimestamp(date string) int64 {
	layout := "01-02-2006"
	t, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("time.Parse has failed: ", err)
	}
	return t.Unix() * 1000
}

// TODO: Insert fuel prices to database.
func InsertFuelPrices(dataList []FuelPrice) {
	//When examining the data from the website, I noticed that the data isn't updated
	//on a daily basis; instead, it is updated whenever new data arrives. Since I want to
	//utilize the daily price changes
	//I'm going to insert the data into the database on a daily basis.
	//I will apply pricing policies for days that are not listed based on the previous data.
	//I will also apply pricing policies for days that are listed but have no data.
	//In the beginning, it will be one time job so I'm not going to afraid of performance issues.
	db := database.Connect()
	defer db.Close()

	// Take each row from the database
	rows, err := db.Query("SELECT * from pricedata")
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var date int64
		var brent float64
		var diesel float64
		err = rows.Scan(&date, &brent, &diesel)
		if err != nil {
			fmt.Println("Scan has failed: ", err)
		}
		for (convertTimestamp(dataList[i].Date) < date) && (convertTimestamp(dataList[i].Date) > date) {
			_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[i].Diesel, date)
			rows.Next()
			if err != nil {
				fmt.Println("Insert has failed: ", err)
			}
		}
		i++
	}
}
