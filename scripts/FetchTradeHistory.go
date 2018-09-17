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
	"os"
)

var o orm.Ormer
func init() {
	parts := []string{beego.AppConfig.String("mysqluser"), ":", beego.AppConfig.String("mysqlpass"),
		"@tcp(", beego.AppConfig.String("mysqlurls"), ":3306)/", beego.AppConfig.String("mysqldb")}
	orm.RegisterDataBase("default", "mysql", strings.Join(parts, ""))
	o = orm.NewOrm()
}

func fetchTradeHistory(stocks []Stock , resolution string, startTime int, endTime int)  {
	for _, stock := range stocks {
		response, _ := http.Get("https://dchart-api.vndirect.com.vn/dchart/history?symbol=" + stock.Symbol + "&resolution=" + resolution +
			"&from=" + strconv.Itoa(startTime) + "&to=" + strconv.Itoa(endTime))
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
			tradeHistory.Resolution = resolution
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
	var (
		stocks []Stock
		startTime, endTime int
	)
	args := os.Args[1:]
	fmt.Println("Start fetching")
	o.QueryTable("stocks").All(&stocks)
	if len(args) > 1 {
		startTime, _ = strconv.Atoi(args[1])
	}
	if len(args) > 2 {
		endTime, _ = strconv.Atoi(args[2])
	}
	fetchTradeHistory(stocks, args[0], startTime, endTime)
}