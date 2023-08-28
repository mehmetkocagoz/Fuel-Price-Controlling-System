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
			_, err := database.Exec(insertQuery, priceList[v].Timestamp, priceList[v].Price)
			if err != nil {
				fmt.Println("Insert has failed: ", err)

			}
		}
	}
}

/*
I don't want to insert same data over and over every time I run the program.
So I'm going to check if there is new data to insert.
*/
func InsertNewBrentOilPrices() {
	// Connect to database
	database := database.Connect()
	defer database.Close()
	// Get the price list
	priceList := GetBrentOilPrices()

	// Get the data from database
	rows, err := database.Query("SELECT timestamp FROM pricedata ORDER BY timestamp ASC")
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	defer rows.Close()
	var isNew bool
	var v int
	for rows.Next() {
		var timestamp int64
		err = rows.Scan(&timestamp)
		if err != nil {
			fmt.Println("Scan has failed: ", err)
		}
		for v = range priceList {
			if timestamp == priceList[v].Timestamp {
				isNew = false
				break
			}
			isNew = true
		}
		if isNew {
			// Insert new data to database
			_, err := database.Exec("INSERT INTO pricedata (timestamp, brentoilprice) VALUES ($1, $2)", priceList[v].Timestamp, priceList[v].Price)
			if err != nil {
				fmt.Println("Insert has failed: ", err)
			}
		}
	}
}
