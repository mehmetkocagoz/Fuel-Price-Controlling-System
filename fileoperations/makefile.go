package fileoperations

import (
	"mehmetkocagz/database"
	"os"
	"strconv"
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
