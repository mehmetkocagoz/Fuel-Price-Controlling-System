package datascrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mehmetkocagz/database"
	"net/http"
)

type Response struct {
	SecurityDetails struct {
		SecID       string `json:"secId"`
		NumDecimals string `json:"numDecimals"`
	} `json:"securityDetails"`
	SeriesData [][]interface{} `json:"seriesData"`
}

type BrentOilPrice struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

func GetBrentOilPrices() []BrentOilPrice {
	url := "https://www.bloomberght.com/piyasa/refdata/brent-petrol"
	priceStruct := []BrentOilPrice{}
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("Get request has failed: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Unmarshal has failed: ", err)
	}
	for v := range data.SeriesData {
		timestamp := int64(data.SeriesData[v][0].(float64))
		price := data.SeriesData[v][1].(float64)

		priceStruct = append(priceStruct, BrentOilPrice{timestamp, price})
		/*
			fmt.Println( "Timestamp:", priceStruct.Timestamp)
			fmt.Println( "Price:", priceStruct.Price)
		*/
	}

	return priceStruct
}

// Insert brent oil prices to database
// This function only used for the first time.
func InsertBrentOilPrices(priceList []BrentOilPrice) {
	database := database.Connect()
	defer database.Close()
	{
		insertQuery := `INSERT INTO pricedata (timestamp, brentoilprice) VALUES ($1, $2)`
		for v := range priceList {
			fmt.Println(v, priceList[v].Timestamp, priceList[v].Price)
			_, err := database.Exec(insertQuery, priceList[v].Timestamp, priceList[v].Price)
			if err != nil {
				fmt.Println("Insert has failed: ", err)

			}
		}
	}
}

/*
I don't want to insert same data over and over every time I run the program.
So I'm going to check the last row in the database and the last row in the priceList.
If the last row in the database is the same as the last row in the priceList
then I'm not going to insert data.
If the last row in the database is not the same as the last row in the priceList
then I'm going to insert the data.
*/
func InsertNewBrentOilPrices() {
	// Connect to database
	database := database.Connect()
	defer database.Close()
	// Get the price list
	priceList := GetBrentOilPrices()
	// Get the last row in the priceList
	lastRowPriceList := priceList[len(priceList)-1]

	// Get the last row in the database
	// I'm going to use this query to get the last row in the database:
	// timestamp DESC LIMIT 1 says that I want to get the last row in the database
	var lastRow BrentOilPrice
	err := database.QueryRow("SELECT timestamp FROM pricedata ORDER BY timestamp DESC LIMIT 1").Scan(&lastRow.Timestamp)
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}

	// Checking if the last row in the database is the same as the last row in the priceList
	if lastRow.Timestamp == lastRowPriceList.Timestamp {
		fmt.Println("There is no new data to insert.")
	} else {
		// If the last row in the database is not the same as the last row in the priceList
		// then I'm going to insert the data.
		// First I'm going to find the index of the lastRow.Timestamp's value in the priceList
		// Then I'm going to insert the data from that index to the end of the priceList
		// to the database.
		var index int
		for i := range priceList {
			if priceList[i].Timestamp == lastRow.Timestamp {
				index = i
				break
			}
		}
		for range priceList[index:] {
			insertQuery := `INSERT INTO pricedata (timestamp, brentoilprice) VALUES ($1, $2)`
			_, err := database.Exec(insertQuery, priceList[index].Timestamp, priceList[index].Price)
			if err != nil {
				fmt.Println("Insert has failed: ", err)
			}
		}
	}
}
