package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type CreateStocksTable_20180906_224742 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &CreateStocksTable_20180906_224742{}
	m.Created = "20180906_224742"

	migration.Register("CreateStocksTable_20180906_224742", m)
}

// Run the migrations
func (m *CreateStocksTable_20180906_224742) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL(`CREATE TABLE stocks (
		id int(11) NOT NULL,
		symbol varchar(255) NULL,
		company_name varchar(255) NULL,
		floor varchar(255) NULL,
		PRIMARY KEY (id) ,
		UNIQUE INDEX symbol (symbol ASC) USING BTREE
	);`)
}

// Reverse the migrations
func (m *CreateStocksTable_20180906_224742) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE stocks")
}
