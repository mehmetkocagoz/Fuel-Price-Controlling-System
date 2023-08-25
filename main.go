package main

import "mehmetkocagz/datascrape"

func main() {

	/*priceList := datascrape.GetBrentOilPrices()
	datascrape.InsertBrentOilPrices(priceList)*/
	datascrape.InsertFuelPrices(datascrape.ScrapeDateAndFuelPrices(*datascrape.GetFuelPrices()))
}
