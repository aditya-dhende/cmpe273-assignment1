package commonService

import (
	"math"
	"strconv"
	"strings"
	"log"
	"fmt"
	"net/http"
)

var tradeId int = 0
var tradeIdWithSharesMap map[int]string = make(map[int]string)
var tradeIdWithUninvestedAmount map[int]float64 = make(map[int]float64)

//code for buy

type ServerStruct struct {}
type Request struct {
    StockSymbolAndPercentage string
    Budget float32
}

type Reply struct {
    TradeId int
    Stocks string
    UnivestedAmount float32
}

func (s *ServerStruct) Buy(httpRq *http.Request,req *Request, res *Reply) error {

	stockCodesList := req.StockSymbolAndPercentage
	totalBudget := req.Budget

	tradeId,stocks,uninvestedamout := getResults(stockCodesList,totalBudget)

	res.TradeId = tradeId
	res.Stocks = stocks
	res.UnivestedAmount = float32(uninvestedamout)

		return nil
}

func getResults(stockCodesList string, totalBudget float32)(int,string,float64) {
	newTradeId := getTradeId()
	var totalAllotedBudget float64 = 0
	var totalRemainingAmount float64 = 0
	var returnStocksString string = ""
	stocksWithPercentage := strings.Split(stockCodesList, ",")

	for data := range stocksWithPercentage {

			singleStockWithPercentage := strings.Split(stocksWithPercentage[data],":")
			percentageWithSymbol := strings.Split(singleStockWithPercentage[1],"%")
			percentageWithoutSymbol,err := strconv.ParseFloat(percentageWithSymbol[0],32)
			if(err!=nil){
				log.Fatal("Unable to parse float in buy function")
			}

			singleStock := singleStockWithPercentage[0]

			singleStockPrice := checkQuote(singleStock)

			allotedBudgetForCurrentStock := float64(totalBudget) * float64(percentageWithoutSymbol) / 100

			stocksPurchasedForCurrentStockFloat := allotedBudgetForCurrentStock / float64(singleStockPrice)
			stocksPurchasedForCurrentStockRounded := math.Floor(float64(stocksPurchasedForCurrentStockFloat))

			leftOverAmount := allotedBudgetForCurrentStock - (stocksPurchasedForCurrentStockRounded * float64(singleStockPrice))
			totalRemainingAmount += leftOverAmount
			totalAllotedBudget += allotedBudgetForCurrentStock

			convertString := singleStock + ":" + strconv.FormatInt(int64(stocksPurchasedForCurrentStockRounded), 10) + ":"+strconv.FormatFloat(float64(singleStockPrice), 'f', 2, 32)

			if returnStocksString==""{
				returnStocksString=convertString
			}else{
				returnStocksString=returnStocksString+","+convertString
			}
	}

				tradeIdWithSharesMap[newTradeId] = returnStocksString
				tradeIdWithUninvestedAmount[newTradeId] = totalRemainingAmount+(float64(totalBudget)-totalAllotedBudget)

				return newTradeId,returnStocksString,totalRemainingAmount+(float64(totalBudget)-totalAllotedBudget)
}


func getTradeId() int{
	tradeId = tradeId +1
	return tradeId
}

//Code for get portfolio
type PortFolioReply struct {
    Stocks string
    CurrentMarketValue float32
    UnivestedAmount float32
    ErrorMessage string
}

type PortFolioRequest struct{
	TradeId int
}

func (s *ServerStruct) Check(httpRq *http.Request, portFolioReq *PortFolioRequest, res *PortFolioReply) error {

	tradeId:=portFolioReq.TradeId
	res.Stocks = ""
	res.CurrentMarketValue = 0
	res.UnivestedAmount=0
	res.ErrorMessage = ""

	var totalMarketValue float64 = 0
	var resultStocks string = ""

	if valShares, ok := tradeIdWithSharesMap[tradeId]; !ok {
		res.ErrorMessage = "Trade Id not found in the database!"
	}else{
		if valUninvestedAmount,ok2 := tradeIdWithUninvestedAmount[tradeId]; ok2{
				stocksWithAmountAndPrice := strings.Split(valShares, ",")
				for data := range stocksWithAmountAndPrice {
					splitStocksWithAmountAndPrice := strings.Split(stocksWithAmountAndPrice[data],":")
					stockCode := splitStocksWithAmountAndPrice[0]
					stockCount := splitStocksWithAmountAndPrice[1]
					stockCountFloat64,err :=  strconv.ParseFloat(stockCount, 64)

					if(err!=nil){
						log.Fatal("Unable to parse float in check function")
					}

					stockOriginalPrice := splitStocksWithAmountAndPrice[2]
					stockOriginalPriceFloat64,err := strconv.ParseFloat(stockOriginalPrice, 64)

										if(err!=nil){
											log.Fatal("Unable to parse float in check function")
										}

					stockOriginalPriceFloatPrecision := fmt.Sprintf("%.2f",stockOriginalPriceFloat64)

					stockPrice := checkQuote(stockCode)
					stockPriceFloat64 := float64(stockPrice)

					var stockCheckedString string = ""


					if stockOriginalPriceFloatPrecision==stockOriginalPriceFloatPrecision{
						stockCheckedString = stocksWithAmountAndPrice[data]
					}else if stockOriginalPriceFloatPrecision > stockOriginalPriceFloatPrecision{
						stockCheckedString = stockCode + ":" + stockCount + ":+" + strconv.FormatFloat(stockPriceFloat64, 'f', 2, 32)
					}else{
						stockCheckedString = stockCode + ":" + stockCount + ":-" + strconv.FormatFloat(stockPriceFloat64, 'f', 2, 32)
					}

					if resultStocks==""{
						resultStocks = stockCheckedString
					}else{
						resultStocks = resultStocks + "," + stockCheckedString
					}

					totalMarketValue += (stockPriceFloat64 * stockCountFloat64)
				}

				res.Stocks = resultStocks
				res.CurrentMarketValue = float32(totalMarketValue)
				res.UnivestedAmount= float32(valUninvestedAmount)
				res.ErrorMessage=""
		}else{
			res.ErrorMessage ="Trade Id not found in the database!"
	}
}

		return nil
}
