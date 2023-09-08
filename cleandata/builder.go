package cleandata

import (
	"fmt"
	"mehmetkocagz/database"
	"mehmetkocagz/datascrape"
	"os"
	"strconv"
	"time"
)

func TableBuilder() {
	database := database.Connect()
	defer database.Close()

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS prices (
		timestamp bigint PRIMARY KEY,
		brentoilprice DECIMAL(10, 2) DEFAULT 0.00,
		fuelprice DECIMAL(10, 2) DEFAULT 0.00,
		exchangeRate DECIMAL(10, 2) DEFAULT 0.00
	);
	`
	_, err := database.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	fmt.Println("Table created successfully.")
}

func FillTableTimestamp() {
	database := database.Connect()
	defer database.Close()

	loc, _ := time.LoadLocation("Europe/Istanbul")
	startDate := time.Date(2018, time.August, 28, 0, 0, 0, 0, loc)
	endDate := time.Date(2023, time.September, 5, 0, 0, 0, 0, loc)

	for startDate.Unix() < endDate.Unix() {
		insertQuery := `
		INSERT INTO prices (timestamp) VALUES ($1)
		`
		_, err := database.Exec(insertQuery, startDate.Unix()*1000)
		if err != nil {
			fmt.Println("Error inserting timestamp:", err)
			return
		}
		startDate = startDate.AddDate(0, 0, 1)
	}
	fmt.Println("Timestamps inserted successfully.")
}

func FillTableBrentOilPrice() {
	// Connect to database
	database := database.Connect()
	defer database.Close()

	// Get brentoilprices from bloomberght.com
	brentoilPrices := datascrape.GetBrentOilPrices()

	// Insert brentoilprices to database
	insertQuery := `
		UPDATE prices SET brentoilprice = $1 WHERE timestamp = $2
		`
	for v := range brentoilPrices {
		_, err := database.Exec(insertQuery, brentoilPrices[v].Price, brentoilPrices[v].Timestamp)
		if err != nil {
			fmt.Println("Error inserting brentoilprice:", err)
			return
		}
	}
}

func FillTableFuelPrice() {
	// Connect to database
	database := database.Connect()
	defer database.Close()

	// Get fuelprices from tppd.com.tr
	doc := *datascrape.GetFuelPrices()
	fuelPriceList := datascrape.ScrapeDateAndFuelPrices(doc)

	// Take each row from the database
	rows, err := database.Query("SELECT * from prices ORDER BY timestamp ASC")
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var date int64
		var brent float64
		var diesel float64
		var exchange float64
		err = rows.Scan(&date, &brent, &diesel, &exchange)
		if err != nil {
			fmt.Println("Scan has failed: ", err)
		}
		if i+1 < len(fuelPriceList) {
			if (fuelPriceList[i].Date <= date) && (fuelPriceList[i+1].Date > date) {
				_, err := database.Exec("UPDATE prices SET fuelprice = $1 WHERE timestamp = $2", fuelPriceList[i].Diesel, date)
				if err != nil {
					fmt.Println("Update has failed: ", err)
				}
			} else {
				i++
				_, err := database.Exec("UPDATE prices SET fuelprice = $1 WHERE timestamp = $2", fuelPriceList[i].Diesel, date)
				if err != nil {
					fmt.Println("Update has failed: ", err)
				}
			}
		} else if i+1 == len(fuelPriceList) {
			_, err := database.Exec("UPDATE prices SET fuelprice = $1 WHERE timestamp = $2", fuelPriceList[i].Diesel, date)
			if err != nil {
				fmt.Println("Insert has failed: ", err)
			}
		}
	}
}

func FillTableExchangeRate() {
	// Connect to database
	database := database.Connect()
	defer database.Close()

	// Get usd exchange rate from bloomberght.com
	usdExchangeRate := datascrape.GetUSDExchangeRate()

	var timestamp int64
	var price float64
	var fuelPrice float64
	var usdExchangeRateColumn float64

	// Take each row from the database
	rows, err := database.Query("SELECT * from prices ORDER BY timestamp ASC")
	if err != nil {
		fmt.Println("Query has failed: ", err)
	}
	defer rows.Close()

	// As I know from API, the data's timestamp is same as brentoilprices.
	// So I can use the timestamp directly.
	for rows.Next() {
		scanError := rows.Scan(&timestamp, &price, &fuelPrice, &usdExchangeRateColumn)
		if scanError != nil {
			fmt.Println("Scan has failed: ", scanError)
		}
		for v := range usdExchangeRate {
			if timestamp == usdExchangeRate[v].Timestamp {
				_, updateError := database.Exec("UPDATE prices SET exchangerate = $1 WHERE timestamp = $2", usdExchangeRate[v].Price, timestamp)
				if updateError != nil {
					fmt.Println("Update has failed: ", updateError)
				}
				break
			}
		}
	}

}
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

	rows, _ := database.Query("SELECT * FROM prices ORDER BY timestamp ASC")
	defer rows.Close()

	file, _ := os.Create("cleanData.csv")
	defer file.Close()

	// Write the header
	writeHeader(file)

	for rows.Next() {
		var timestamp int64
		var price float64
		var fuelPrice float64
		var usdExchangeRateColumn float64
		scanError := rows.Scan(&timestamp, &price, &fuelPrice, &usdExchangeRateColumn)
		if scanError != nil {
			panic(scanError)
		}
		file.WriteString(strconv.FormatFloat(price, 'f', -1, 64) + "," + strconv.FormatFloat(fuelPrice, 'f', -1, 64) + "," + strconv.FormatFloat(usdExchangeRateColumn, 'f', -1, 64) + "\n")
	}
}
