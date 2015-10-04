package commonService

import (
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"log"
	"strconv"
	"net/http"
)


func checkQuote(stockName string) float32 {
	baseUrlLeft := "https://query.yahooapis.com/v1/public/yql?q=select%20LastTradePriceOnly%20from%20yahoo.finance%0A.quotes%20where%20symbol%20%3D%20%22"
	baseUrlRight := "%22%0A%09%09&format=json&env=http%3A%2F%2Fdatatables.org%2Falltables.env"

	resp, err := http.Get(baseUrlLeft + stockName + baseUrlRight)

	if err != nil {
		log.Fatal("Error on http request for quote")
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		log.Fatal("Error reading response body for json")
	}

	if resp.StatusCode != 200 {
		log.Fatal("Query failure, please check network connection and stock code for issues ")
	}

	newSimpleJSON, err := simplejson.NewJson(body)
	if err != nil {
		log.Fatal("Error getting JSON from body")
	}

	price, _ := newSimpleJSON.Get("query").Get("results").Get("quote").Get("LastTradePriceOnly").String()
	floatPrice, err := strconv.ParseFloat(price, 32)

	return float32(floatPrice)
}
