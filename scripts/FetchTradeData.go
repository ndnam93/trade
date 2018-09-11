package main

import (
	"github.com/astaxie/beego/orm"
	"fmt"
	"strings"
	"github.com/astaxie/beego"
	"net/http"
	"io/ioutil"
	"github.com/bitly/go-simplejson"
	"trade/models"
	"os"
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

func main() {
	symbolsResponse, _ := http.Get("https://finfo-api.vndirect.com.vn//stocks/mini")
	res, _ := ioutil.ReadAll(symbolsResponse.Body)

	json, _ := simplejson.NewJson(res)
	list, _ := json.Get("data").Array()


	o := orm.NewOrm()
	for _, c := range list {
		company, _ := c.(map[interface{}]interface{})
		fmt.Println(company)
		os.Exit(0)
		var stock models.Stock
		stock.Symbol = company["symbol"].(string)
		fmt.Println(stock)
		os.Exit(0)
		o.Insert(&stock)
	}

}