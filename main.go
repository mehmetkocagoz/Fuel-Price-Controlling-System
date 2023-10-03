package main

import (
	"fmt"
	"mehmetkocagz/cleandata"
	"mehmetkocagz/datascrape"
	"mehmetkocagz/fileoperations"
	"mehmetkocagz/handlers"
	"net/http"
	"os/exec"
	"time"

	"github.com/gorilla/mux"
)

// This function will fill the database with the data we want.
// This function will be used only once.
// After that we will just update the database if new data comes.
func databaseFiller() {
	// Get brent oil prices from bloomberght.com
	priceList := datascrape.GetBrentOilPrices()
	// Insert brent oil prices to database
	datascrape.InsertBrentOilPrices(priceList)
	// Get fuel prices from tppd.com.tr
	doc := datascrape.GetFuelPrices()
	// Insert fuel prices to database
	datascrape.InsertFuelPrices(datascrape.ScrapeDateAndFuelPrices(*doc))
	// Update usd exchange rate
	datascrape.UpdateUSDExchangeRate()
	// Get csv file from database
	fileoperations.CreateAndWritetoCSV()
}
func databaseUpdater() {
	// First insert new brent oil prices
	datascrape.InsertNewBrentOilPrices()
	// Then insert new fuel prices
	datascrape.InsertNewFuelPrices()
	// Update usd exchange rate
	datascrape.UpdateUSDExchangeRate()
	// Update csv file
	fileoperations.UpdateCSVFile()
}

func cleanedDataFiller() {
	cleandata.TableBuilder()
	cleandata.FillTableTimestamp()
	cleandata.FillTableBrentOilPrice()
	cleandata.FillTableFuelPrice()
	cleandata.FillTableExchangeRate()
	cleandata.CreateAndWritetoCSV()
}

func cleanedDataUpdater() {
	cleandata.UpdateTableTimestamp()
	cleandata.UpdateTableBrentOilPrice()
	cleandata.UpdateTableFuelPrice()
	cleandata.UpdateTableExchangeRate()
	cleandata.UpdateCSV()
}

func linearRegression() {
	cmd := exec.Command("python", "datafunctions/linearRegression.py")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
		return
	}
	fmt.Println("A", string(out))
	fmt.Println("Linear regression script executed successfully.")
}

func runServer() {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./assests/"))
	r.PathPrefix("/assests").Handler(http.StripPrefix("/assests", fs))

	// Home page
	homeRouter := r.Methods("GET").Subrouter()
	homeRouter.HandleFunc("/", handlers.ServeHome)

	// Analysis Chart page
	analysisRouter := r.Methods("GET").Subrouter()
	analysisRouter.HandleFunc("/analytic.html", handlers.ServeAnalysis)

	// Create a new server
	srv := &http.Server{
		Addr:         ":9090",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server is running on port 9090.")
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func main() {
	//databaseUpdater()
	//runServer()

	cleandata.UpdateTableTimestamp()
	cleandata.UpdateTableBrentOilPrice()
	cleanedDataFiller()
}
