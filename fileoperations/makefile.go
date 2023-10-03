package fileoperations

import (
	"encoding/json"
	"fmt"
	"log"
	"mehmetkocagz/database"
	"mehmetkocagz/model"
	"os"
	"strconv"
	"time"
)

func writeHeader(file *os.File) {
	header := []string{"BrentOilPrice", "FuelPrice", "USD/TRY"}
	for _, v := range header {
		if v == "USD/TRY" {
			file.WriteString(v + "\n")

		} else {
			file.WriteString(v + ",")
		}
	}
}

func CreateAndWritetoCSV() {
	// Connect to database
	database := database.Connect()

	rows, _ := database.Query("SELECT * FROM pricedata ORDER BY timestamp ASC")
	defer rows.Close()

	file, _ := os.Create("data.csv")
	defer file.Close()

	// Write the header
	writeHeader(file)

	for rows.Next() {
		var timestamp int64
		var price float64
		var fuelPrice float64
		var dateColumn string
		var usdExchangeRateColumn float64
		scanError := rows.Scan(&timestamp, &price, &fuelPrice, &dateColumn, &usdExchangeRateColumn)
		if scanError != nil {
			panic(scanError)
		}
		file.WriteString(dateColumn + "," + strconv.FormatFloat(price, 'f', -1, 64) + "," + strconv.FormatFloat(fuelPrice, 'f', -1, 64) + "," + strconv.FormatFloat(usdExchangeRateColumn, 'f', -1, 64) + "\n")
	}
}

func UpdateCSVFile() {
	//Connect to database
	database := database.Connect()

	rows, _ := database.Query("SELECT * FROM pricedata ORDER BY timestamp ASC")
	defer rows.Close()

	file, _ := os.OpenFile("data.csv", os.O_RDWR, 0666)
	defer file.Close()

	// Let's delete the old data in the file
	file.Truncate(0)

	// Write the header
	writeHeader(file)

	for rows.Next() {
		var timestamp int64
		var price float64
		var fuelPrice float64
		var dateColumn string
		var usdExchangeRateColumn float64
		scanError := rows.Scan(&timestamp, &price, &fuelPrice, &dateColumn, &usdExchangeRateColumn)
		if scanError != nil {
			panic(scanError)
		}
		file.WriteString(strconv.FormatFloat(price, 'f', -1, 64) + "," + strconv.FormatFloat(fuelPrice, 'f', -1, 64) + "," + strconv.FormatFloat(usdExchangeRateColumn, 'f', -1, 64) + "\n")
	}
}

func ConvertJSON() {
	db := database.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM prices order by timestamp")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var data []model.PriceAnalysis
	var hold model.PriceAnalysis
	var timestamp int64
	for rows.Next() {

		err := rows.Scan(&timestamp, &hold.BrentPrice, &hold.FuelPrice, &hold.FuelPrice)
		if err != nil {
			log.Fatal(err)
		}

		hold.PriceAnalysisID = time.Unix(timestamp/1000, 0)
		fmt.Printf("Timestamp: %s\n", hold.PriceAnalysisID)
		data = append(data, hold)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	outputFile, err := os.Create("output.json")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	_, err = outputFile.Write(jsonData)
	if err != nil {
		log.Fatal(err)
	}
}
