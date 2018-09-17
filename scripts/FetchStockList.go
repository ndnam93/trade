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

var o orm.Ormer
func init() {
	parts := []string{beego.AppConfig.String("mysqluser"), ":", beego.AppConfig.String("mysqlpass"),
		"@tcp(", beego.AppConfig.String("mysqlurls"), ":3306)/", beego.AppConfig.String("mysqldb")}
	orm.RegisterDataBase("default", "mysql", strings.Join(parts, ""))
	o = orm.NewOrm()
}

type Company struct {
	id int
	symbol string
	company string
	floor string
	indexCode string
}

type Companies []Company


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


func main() {
	fmt.Println("Start fetching")
	fetchStockList()
}