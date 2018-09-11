package main

import (
	_ "trade/routers"
	"github.com/astaxie/beego"
	"os"
	"github.com/astaxie/beego/orm"
	"strings"
)

func init() {
	parts := []string{os.Getenv("MYSQL_USER"), ":", os.Getenv("MYSQL_PASSWORD"),
		"@tcp(", os.Getenv("MYSQL_HOST"), ":3306)/", os.Getenv("MYSQL_DATABASE")}
	orm.RegisterDataBase("default", "mysql", strings.Join(parts, ""))
}

func main() {
	beego.Run()
}

