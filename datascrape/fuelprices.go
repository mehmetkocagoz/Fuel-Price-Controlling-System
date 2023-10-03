package datascrape

import (
	"fmt"
	"mehmetkocagz/database"
	"mehmetkocagz/datafunctions"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FuelPrice struct {
	Date   int64
	Diesel float64
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
	var date int64
	var diesel float64
	doc.Find("table tr").Each(func(a int, s *goquery.Selection) {
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			// As I know how data table structered I can get the data I want.
			// I'm going to get the date and DIESEL prices only.
			// This function can be improved.
			if i == 0 {
				// Date will come to us as string. 13 May 2019
				// I'm going to convert it to int64.
				date = datafunctions.ConvertTimestamp(s.Text())
			} else if i == 4 {
				// Assume there is no error.
				// We can handle it later.
				// s.Text() has white spaces at the end of the string.
				// I'm going to trim it.
				dieselString := strings.TrimSpace(s.Text())
				// I'm going to convert it to float64.
				diesel, _ = strconv.ParseFloat(dieselString, 64)

			}
		})

		fuelPrices = append(fuelPrices, FuelPrice{date, diesel})
	})
	// fuelPrices has the data from 2008
	// I'm going to remove the data before August 2018.
	// I'm going to use the data from August 2018.
	// 1534550400000 18 August 2018
	indexOfStartingDate := 0
	for range fuelPrices {
		if fuelPrices[indexOfStartingDate].Date > 1534550400000 {
			fuelPrices = fuelPrices[indexOfStartingDate-1:]
			fmt.Println("Success")
			break
		} else {
			indexOfStartingDate++
		}
	}
	return fuelPrices
}

// TODO: Update fuel prices to database.
// pricedata table has fuelprice column.
// In default, fuelprice column is 0.
// I'm going to update the fuelprice column with the data I get from tppd.com.tr
// Maybe later, I can add a new function that updates fuelprices only if the data is new.
// But now, This function will update table's every row every time I run the program.
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
	rows, err := db.Query("SELECT * from pricedata ORDER BY timestamp ASC")
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	defer rows.Close()
	i := 0
	for rows.Next() {
		var date int64
		var brent float64
		var diesel float64
		var dateStr string
		var exchange float64
		err = rows.Scan(&date, &brent, &diesel, &dateStr, &exchange)
		if err != nil {
			fmt.Println("Scan has failed: ", err)
		}
		if i+1 < len(dataList) {
			if (dataList[i].Date <= date) && (dataList[i+1].Date > date) {
				_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[i].Diesel, date)
				if err != nil {
					fmt.Println("Update has failed: ", err)
				}
			} else {
				i++
				_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[i].Diesel, date)
				if err != nil {
					fmt.Println("Update has failed: ", err)
				}
			}
		} else if i+1 == len(dataList) {
			_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[i].Diesel, date)
			if err != nil {
				fmt.Println("Insert has failed: ", err)
			}
		}
	}
}
func InsertNewFuelPrices() {
	// Connect to database
	db := database.Connect()
	defer db.Close()
	// Get the price list
	dataList := ScrapeDateAndFuelPrices(*GetFuelPrices())
	// Count the number of rows that has fuelprice = 0
	var count int
	err := db.QueryRow("SELECT COUNT(*) from pricedata WHERE fuelprice=0").Scan(&count)
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	// Take each row from the database
	rows, err := db.Query("SELECT * from pricedata WHERE fuelprice=0 ORDER BY timestamp ASC")
	for rows.Next() {
		var timestamp int64
		var brent float64
		var diesel float64
		var dateStr string
		var exchange float64
		err = rows.Scan(&timestamp, &brent, &diesel, &dateStr, &exchange)
		if err != nil {
			fmt.Println("Scan has failed: ", err)
		}
		// If there is only one row that has fuelprice = 0
		// I'm going to update it with the last data I get from tppd.com.tr
		if count == 1 {
			_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[len(dataList)-1].Diesel, timestamp)
			if err != nil {
				fmt.Println("Update has failed: ", err)
			}
			return
		} else {
			// If there is more than one row that has fuelprice = 0
			// First I'm going to slice the dataList to get the data I want.
			// With slicing, Function will work faster because there will be less data to check.
			// Then I'm going to update the data with the data I get from tppd.com.tr
			dataList = dataList[len(dataList)-count:]
			for v := range dataList {
				if timestamp >= dataList[v].Date && timestamp < dataList[v+1].Date {
					_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[v].Diesel, timestamp)
					if err != nil {
						fmt.Println("Update has failed: ", err)
					}
					count--
					break
				}
			}
		}
	}
}
