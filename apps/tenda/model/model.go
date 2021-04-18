package model

import (
	"fmt"
	"database/sql"
	"github.com/guangxue/webpages/mysql"
	_ "github.com/go-sql-driver/mysql"	
)
var db = mysql.Connect("tenda");

func GetAll() *sql.Rows {
	rows, err := db.Query("select distinct model from stock_locations")
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	return rows
}

func Close() {
	db.Close()
}