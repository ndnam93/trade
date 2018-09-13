package models

import "github.com/astaxie/beego/orm"

type TradeHistory struct {
	Time int `orm:"pk""`
	Symbol string
	Close int
	Open int
	High int
	Low int
	Volume int
}

func (tradeHistory *TradeHistory) TableName() string {
	return "trade_history"
}

func init() {
	orm.RegisterModel(new(TradeHistory))
}