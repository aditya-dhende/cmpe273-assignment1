package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/bitly/go-simplejson"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"bufio"
	"os"
	"strings"
)


func main() {

	var choiceInt int = 0
	reader := bufio.NewReader(os.Stdin)

	for !(choiceInt==3){
	fmt.Println("\n\nMenu : ")
	fmt.Println("1. Buy Stocks")
	fmt.Println("2. Check portfolio with trade ID")
	fmt.Println("3. Exit")
	fmt.Print("Enter your choice number : ")
	choice, err := reader.ReadString('\n')
	if(err!= nil){
	  log.Fatal("error while reading choice")
	}
	choice = strings.Replace(choice,"\r","", -1)
	choice = strings.Replace(choice,"\n","", -1)
	fmt.Println("You have selected option: "+choice)
	choiceInt , err := strconv.Atoi(choice)

	if err != nil {
	  log.Fatal("Int conversion error for choice. Did you enter an invalid choice?")
	}


	switch choiceInt{

	case 1 :
	fmt.Print("\nEnter the stock symbol and percentage string. E.g GOOG:50%,YHOO:50% :  ")
	stockSymPercent, _ := reader.ReadString('\n')
	stockSymPercent = strings.Replace(stockSymPercent,"\r","", -1)
	stockSymPercent = strings.Replace(stockSymPercent,"\n","", -1)

	var percentageTotal float64 =0
	stocksWithPercentage := strings.Split(stockSymPercent, ",")
	for data := range stocksWithPercentage{

		singleStockWithPercentage := strings.Split(stocksWithPercentage[data],":")
		percentageWithSymbol := strings.Split(singleStockWithPercentage[1],"%")
		percentageWithoutSymbol,err := strconv.ParseFloat(percentageWithSymbol[0],32)
		if(err!=nil){
			log.Fatal("Unable to parse percentage amount to float")
		}
		percentageTotal += percentageWithoutSymbol
	}
	if (percentageTotal>100)||(percentageTotal<0){
		fmt.Println("Percentage Total is invalid!")
	}else{
	fmt.Print("Enter the budget :  ")
	budget, _ := reader.ReadString('\n')
	budget = strings.Replace(budget,"\r","", -1)
	budget = strings.Replace(budget,"\n","", -1)
	budgetFloat,err := strconv.ParseFloat(budget, 32)

	if err != nil {
		log.Fatal("Float conversion error for budget. Did you enter a number? ")
	}

	data, err := json.Marshal(map[string]interface{}{
				"method": "ServerStruct.Buy",
				"id":     1,
				"params": []map[string]interface{}{map[string]interface{}{"StockSymbolAndPercentage": stockSymPercent, "Budget": float32(budgetFloat)}},
			})

		if err != nil {
			log.Fatal("Error while marshalling Buy request ")
		}

		resp, err := http.Post("http://127.0.0.1:9999/", "application/json", strings.NewReader(string(data)))

		if err != nil {
			log.Fatal("Error while posting data to the server for buy request")
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal("Error while reading response body for buy request")
		}

		newSimplejson, err := simplejson.NewJson(body)

		if err != nil {
			log.Fatal("Error while unmarshalling json from response body for buy request")
		}


		tradeId, _ := newSimplejson.Get("result").Get("TradeId").Int()
		stocks,_ := newSimplejson.Get("result").Get("Stocks").String()
	  uninvestedAmount, _ := newSimplejson.Get("result").Get("UnivestedAmount").Float64()


		fmt.Print("\nTrade Id: ")
		fmt.Println(tradeId)
		fmt.Println("Stocks Purchased: " + stocks)
		fmt.Println("Uninvested Amount: "+strconv.FormatFloat(uninvestedAmount, 'f', 2, 32))
}

case 2 :
  fmt.Print("\nEnter the trade ID :  ")
  tradeId, _ := reader.ReadString('\n')

	tradeId = strings.Replace(tradeId,"\r","", -1)
	tradeId = strings.Replace(tradeId,"\n","", -1)

  tradeIdInt , err := strconv.Atoi(tradeId)

  if err != nil {
  	log.Fatal("Int conversion error while reading trade id. Did you enter an integer value?")
  }

		data, err := json.Marshal(map[string]interface{}{
					"method": "ServerStruct.Check",
					"id":     1,
					"params": []map[string]interface{}{map[string]interface{}{"TradeId": tradeIdInt}},
				})

	if err != nil {
		log.Fatal("Error while marshalling Check request ")
	}

	resp, err := http.Post("http://127.0.0.1:9999/", "application/json", strings.NewReader(string(data)))

	if err != nil {
		log.Fatal("Error while posting data to the server for check request")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("Error while reading response body for check request")
	}

	newSimplejson, err := simplejson.NewJson(body)

	if err != nil {
		log.Fatal("Error while unmarshalling json from response body for check request")
	}

			stocks, _ := newSimplejson.Get("result").Get("Stocks").String()

			uninvestedAmount, _ := newSimplejson.Get("result").Get("UnivestedAmount").Float64()

			totalMarketValue, _ := newSimplejson.Get("result").Get("CurrentMarketValue").Float64()

			errorMessage, _ := newSimplejson.Get("result").Get("ErrorMessage").String()

	if errorMessage==""{
		fmt.Println("\nStocks: "+stocks)
		fmt.Println("Uninvested amount: " + strconv.FormatFloat(uninvestedAmount, 'f', 2, 32))
		fmt.Println("Current Market Value: "+strconv.FormatFloat(totalMarketValue, 'f', 2, 32))
	}else{
		fmt.Println(errorMessage)
	}


case 3: fmt.Println("Exiting")
return

default: fmt.Println("Invalid Choice")
}

}
}
