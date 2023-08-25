package datascrape

import (
	"fmt"
	"mehmetkocagz/database"
	"net/http"
	"strconv"
	"strings"
	"time"

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
				date = convertTimestamp(s.Text())
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

func switchMonthToNumber(month string) string {
	switch month {
	case "January":
		return "01"
	case "February":
		return "02"
	case "March":
		return "03"
	case "April":
		return "04"
	case "May":
		return "05"
	case "June":
		return "06"
	case "July":
		return "07"
	case "August":
		return "08"
	case "September":
		return "09"
	case "October":
		return "10"
	case "November":
		return "11"
	case "December":
		return "12"
	}
	return "0"
}

func convertTimestamp(date string) int64 {
	//fmt.Println("converting..", date)
	// I know that our date will come like int string int format.
	// So first I'm going to convert it to int-int-int format.
	parsedDate := strings.Split(date, " ")
	month := switchMonthToNumber(parsedDate[1])
	date = parsedDate[2] + "-" + month + "-" + parsedDate[0]
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)

	if err != nil {
		fmt.Println("time.Parse has failed: ", err)
	}
	return (t.Unix() * 1000)
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
		err = rows.Scan(&date, &brent, &diesel)
		if err != nil {
			fmt.Println("Scan has failed: ", err)
		}
		if i+1 < len(dataList) {
			if (dataList[i].Date <= date) && (dataList[i+1].Date > date) {
				_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[i].Diesel, date)
				if err != nil {
					fmt.Println("Insert has failed: ", err)
				}
			} else {
				i++
				_, err := db.Exec("UPDATE pricedata SET fuelprice = $1 WHERE timestamp = $2", dataList[i].Diesel, date)
				if err != nil {
					fmt.Println("Insert has failed: ", err)
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
