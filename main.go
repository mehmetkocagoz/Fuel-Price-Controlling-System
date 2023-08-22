package main

import (
	"fmt"
	"mehmetkocagz/datascrape"
)

func main() {
	price := datascrape.GetBrentOilPrices()
	fmt.Println(price)
}
