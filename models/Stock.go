package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Stock struct {
	Id   int
	Symbol string `orm:"size(255)"`
	CompanyName string `orm:"size(255)"`
	Floor string `orm:"size(255)"`
}

func (stock *Stock) TableName() string {
	return "stocks"
}

func init() {
	orm.RegisterModel(new(Stock))

	//orm.RegisterDataBase("default", "mysql", "root:@/trade")
}