package main

import (
	"github.com/astaxie/beego/orm"
	"strings"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	. "trade/models"
	"fmt"
	"strconv"
)

func init() {
	parts := []string{beego.AppConfig.String("mysqluser"), ":", beego.AppConfig.String("mysqlpass"),
		"@tcp(", beego.AppConfig.String("mysqlurls"), ":3306)/", beego.AppConfig.String("mysqldb")}
	orm.RegisterDataBase("default", "mysql", strings.Join(parts, ""))
}

type Company struct {
	id int
	symbol string
	company string
	floor string
	indexCode string
}

type Companies []Company

var o = orm.NewOrm()

func fetchStockList() []Stock {
	var (
		createdCount int64
		stocks []Stock
	)
	symbolsResponse, _ := http.Get("https://finfo-api.vndirect.com.vn/stocks/mini")
	res, _ := ioutil.ReadAll(symbolsResponse.Body)

	list := gjson.Get(string(res), "data")
	list.ForEach(func(key, value gjson.Result) bool {
		valueMap := value.Map()
		var stock Stock
		stock.Symbol = valueMap["symbol"].String()
		created, _, _ := o.ReadOrCreate(&stock, "Symbol")
		stock.CompanyName = valueMap["company"].String()
		stock.Floor = valueMap["floor"].String()
		o.Update(&stock)
		if created {
			createdCount += 1
		}
		stocks = append(stocks, stock)
		return true
	})
	fmt.Println("Created " + strconv.Itoa(int(createdCount)) + " and updated " + strconv.Itoa(int(list.Get("#").Int() - createdCount)) + " stocks")
	return stocks
}

func fetchTradeHistory(stocks []Stock , resolution string)  {
	for _, stock := range stocks {
		response, _ := http.Get("https://dchart-api.vndirect.com.vn/dchart/history?symbol=" + stock.Symbol + "&resolution=D")
		res, _ := ioutil.ReadAll(response.Body)
		data := gjson.Parse(string(res))

		var startIndex int
		lastTimestamp := getLastTimestamp(stock.Symbol)
		for i, timestamp := range data.Get("t").Array() {
			if int(timestamp.Int()) > lastTimestamp {
				startIndex = i
			}
		}

		var tradeHistoryList []TradeHistory
		for i := startIndex; i < int(data.Get("t.#").Int()); i++ {
			var tradeHistory TradeHistory
			tradeHistory.Symbol = stock.Symbol
			tradeHistory.Time = int(data.Get("t." + strconv.Itoa(i)).Int())
			tradeHistory.Open = int(data.Get("o." + strconv.Itoa(i)).Int())
			tradeHistory.Close = int(data.Get("c." + strconv.Itoa(i)).Int())
			tradeHistory.High = int(data.Get("h." + strconv.Itoa(i)).Int())
			tradeHistory.Low = int(data.Get("l." + strconv.Itoa(i)).Int())
			tradeHistory.Volume = int(data.Get("v." + strconv.Itoa(i)).Int())
			tradeHistoryList = append(tradeHistoryList, tradeHistory)
		}
		o.InsertMulti(len(tradeHistoryList), tradeHistoryList)

		fmt.Println("Inserted " + strconv.Itoa(len(tradeHistoryList)) + " trade history records for symbol " + stock.Symbol)
	}
}
func getLastTimestamp(symbol string) int {
	var tradeHistory TradeHistory
	error := o.QueryTable("trade_history").
		Filter("symbol", symbol).
		OrderBy("-time").
		Limit(1).
		One(&tradeHistory)
	if error == nil {
		return tradeHistory.Time
	}
	return 0
}


func main() {
	fmt.Println("Start fetching")
	stocks := fetchStockList()
	fetchTradeHistory(stocks, "15")
}