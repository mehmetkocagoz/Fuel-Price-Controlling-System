package main

import "mehmetkocagz/datascrape"

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
}

func main() {
}
