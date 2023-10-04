package model

import (
	"mehmetkocagz/database"
	"time"
)

type PriceAnalysis struct {
	Timestamp    int64   `json:"timestamp"`
	BrentPrice   float64 `json:"brentprice"`
	FuelPrice    float64 `json:"price"`
	ExchangeRate float64 `json:"exchangeRate"`
}

type PriceAnalysisTable struct {
	PriceAnalysisID time.Time `json:"priceAnalysisID"`
	BrentPrice      float64   `json:"brentprice"`
	FuelPrice       float64   `json:"price"`
	ExchangeRate    float64   `json:"exchangeRate"`
}

type PriceAnalysisListData struct {
	LastBrentPrice   float64
	LastFuelPrice    float64
	LastExchangeRate float64
	TableData        []PriceAnalysisTable
	BrentPriceRate   float64
	FuelPriceRate    float64
	ExchangeRateRate float64
}

func GrabTemplateData() PriceAnalysisListData {
	data := grabData()
	var priceAnalysisListData PriceAnalysisListData
	priceAnalysisListData.LastBrentPrice = data[len(data)-1].BrentPrice
	priceAnalysisListData.LastFuelPrice = data[len(data)-1].FuelPrice
	priceAnalysisListData.LastExchangeRate = data[len(data)-1].ExchangeRate
	priceAnalysisListData.TableData = ConvertToPriceAnalysisTable(data[len(data)-6 : len(data)-1])
	priceAnalysisListData.BrentPriceRate = (data[len(data)-1].BrentPrice - data[len(data)-2].BrentPrice) / data[len(data)-2].BrentPrice * 100
	priceAnalysisListData.FuelPriceRate = (data[len(data)-1].FuelPrice - data[len(data)-2].FuelPrice) / data[len(data)-2].FuelPrice * 100
	priceAnalysisListData.ExchangeRateRate = (data[len(data)-1].ExchangeRate - data[len(data)-2].ExchangeRate) / data[len(data)-2].ExchangeRate * 100
	return priceAnalysisListData
}

func ConvertToPriceAnalysisTable(paSlice []PriceAnalysis) []PriceAnalysisTable {
	patSlice := make([]PriceAnalysisTable, len(paSlice))
	i := len(paSlice)

	for _, pa := range paSlice {
		pat := PriceAnalysisTable{
			PriceAnalysisID: time.Unix(pa.Timestamp/1000, 0),
			BrentPrice:      pa.BrentPrice,
			FuelPrice:       pa.FuelPrice,
			ExchangeRate:    pa.ExchangeRate,
		}
		patSlice[i-1] = pat
		i--
	}

	return patSlice
}
cdd
func grabData() []PriceAnalysis {
	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM prices order by timestamp ")
	if err != nil {
		panic(err.Error())
	}
	var priceAnalysisList []PriceAnalysis
	for rows.Next() {
		var priceAnalysis PriceAnalysis
		err = rows.Scan(&priceAnalysis.Timestamp, &priceAnalysis.BrentPrice, &priceAnalysis.FuelPrice, &priceAnalysis.ExchangeRate)
		if err != nil {
			panic(err.Error())
		}
		priceAnalysisList = append(priceAnalysisList, priceAnalysis)
	}
	return priceAnalysisList
}
