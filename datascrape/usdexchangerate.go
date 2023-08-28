package datascrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mehmetkocagz/database"
	"net/http"
)

// Like fuelprices, I'm going to create a struct for usd exchange rate.

type USDExchangeRate struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

// In brentoilprices I used bloomberght.com for data.
// I'm going to use bloomberght.com for usd exchange rate.
// Since both of them are bloomberght.com, I'm going to use the same response struct.
var data Response

// Get usd exchange rate from bloomberght.com
// This function will return usd exchange rate as a struct.
func GetUSDExchangeRate() []USDExchangeRate {

	url := "https://www.bloomberght.com/piyasa/refdata/dolar"

	exchangeStruct := []USDExchangeRate{}
	resp, responseError := http.Get(url)
	if responseError != nil {
		fmt.Println("Get request has failed: ", responseError)
	}
	defer resp.Body.Close()

	body, readingError := ioutil.ReadAll(resp.Body)
	if readingError != nil {
		fmt.Println("Reading body has failed: ", readingError)
	}
	unMarshalError := json.Unmarshal(body, &data)
	if unMarshalError != nil {
		fmt.Println("Unmarshal has failed: ", unMarshalError)
	}
	for v := range data.SeriesData {
		timestamp := int64(data.SeriesData[v][0].(float64))
		price := data.SeriesData[v][1].(float64)
		exchangeStruct = append(exchangeStruct, USDExchangeRate{timestamp, price})
	}

	return exchangeStruct
}

// Like fuelprices it will be one time job to update usd exchange rate.
// I'm going to create a function to update usd exchange rate.
func UpdateUSDExchangeRate() {
	// Connect to database
	database := database.Connect()
	// Get usd exchange rate from bloomberght.com
	usdExchangeRate := GetUSDExchangeRate()

	var timestamp int64
	var price float64
	var fuelPrice float64
	var dateColumn string
	var usdExchangeRateColumn float64

	// Take each row from the database
	rows, err := database.Query("SELECT * from pricedata ORDER BY timestamp ASC")
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	defer rows.Close()

	// As I know from API, the data's timestamp is same as brentoilprices.
	// So I can use the timestamp directly.
	for rows.Next() {
		scanError := rows.Scan(&timestamp, &price, &fuelPrice, &dateColumn, &usdExchangeRateColumn)
		if scanError != nil {
			fmt.Println("Scan has failed: ", scanError)
		}
		for v := range usdExchangeRate {
			if timestamp == usdExchangeRate[v].Timestamp {
				_, updateError := database.Exec("UPDATE pricedata SET exchange_column = $1 WHERE timestamp = $2", usdExchangeRate[v].Price, timestamp)
				if updateError != nil {
					fmt.Println("Update has failed: ", updateError)
				}
				break
			}
		}
	}
}
