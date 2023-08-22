package datascrape

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
