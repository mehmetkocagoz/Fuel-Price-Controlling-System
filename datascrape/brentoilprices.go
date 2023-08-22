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
func InsertBrentOilPrices(priceList []BrentOilPrice) {
	database := database.Connect()
	defer database.Close()
	/*I don't want to insert same data over and over every time I run the program.
		So I'm going to check the last row in the database and the last row in the priceList.
	    If the last row in the database is the same as the last row in the priceList
		then I'm not going to insert anything.
		If the last row in the database is not the same as the last row in the priceList
		then I'm going to insert the data.
	*/

	// Get the last row in the database
	var lastRow BrentOilPrice
	err := database.QueryRow("SELECT timestamp, price FROM brentoil ORDER BY timestamp DESC LIMIT 1").Scan(&lastRow.Timestamp, &lastRow.Price)
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	// Get the last row in the priceList
	lastRowPriceList := priceList[len(priceList)-1]
	// Checking if the last row in the database is the same as the last row in the priceList
	if lastRow.Timestamp == lastRowPriceList.Timestamp {
		fmt.Println("There is no new data to insert.")
	} else {
		insertQuery := `INSERT INTO brentoil (timestamp, price) VALUES ($1, $2)`
		for v := range priceList {
			_, err := database.Exec(insertQuery, priceList[v].Timestamp, priceList[v].Price)
			if err != nil {
				fmt.Println("Insert has failed: ", err)

			}
		}
	}
}
