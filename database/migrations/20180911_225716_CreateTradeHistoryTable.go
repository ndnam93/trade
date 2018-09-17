package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateTradeHistoryTable_20180911_225716 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateTradeHistoryTable_20180911_225716{}
	m.Created = "20180911_225716"

	migration.Register("CreateTradeHistoryTable_20180911_225716", m)
}

// Run the migrations
func (m *CreateTradeHistoryTable_20180911_225716) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE trade_history (
		time int NOT NULL,
		symbol varchar(255) NULL,
		close int(255) NULL,
		open int(255) NULL,
		high int(255) NULL,
		low int(255) NULL,
		volume int(255) NULL,
		resolution ENUM('15', '30', '60', 'D')
		PRIMARY KEY (time)
	);`)
}

// Reverse the migrations
func (m *CreateTradeHistoryTable_20180911_225716) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE trade_history")
}
