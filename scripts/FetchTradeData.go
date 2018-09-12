package main

import (
	"github.com/astaxie/beego/orm"
	"strings"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"trade/models"
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

func fetchStockList() {
	var createdCount int64
	symbolsResponse, _ := http.Get("https://finfo-api.vndirect.com.vn//stocks/mini")
	res, _ := ioutil.ReadAll(symbolsResponse.Body)

	o := orm.NewOrm()
	list := gjson.Get(string(res), "data")
	//fmt.Printf("%#v", strconv.Itoa(int(list.Get("#").Int())))
	//os.Exit(0)
	list.ForEach(func(key, value gjson.Result) bool {
		valueMap := value.Map()
		var stock models.Stock
		stock.Symbol = valueMap["symbol"].String()
		created, _, _ := o.ReadOrCreate(&stock, "Symbol")
		stock.CompanyName = valueMap["company"].String()
		stock.Floor = valueMap["floor"].String()
		o.Update(&stock)
		if created {
			createdCount += 1
		}
		return true
	})
	fmt.Println("Created " + strconv.Itoa(int(createdCount)) + " and updated " + strconv.Itoa(int(list.Get("#").Int() - createdCount)) + " stocks")
}

func main() {
	fetchStockList()
}