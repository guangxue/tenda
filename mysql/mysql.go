package mysql

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"	
)

func Connect(dbname string) *sql.DB {
	dsn := "gzhang:guangxue@/" + dbname
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("DB connection error: ", err)
	}
	var pingerr = db.Ping()
	if pingerr != nil {
		fmt.Println("DB pinging error: ", pingerr)
	}
	fmt.Println("[mysql.Connect] DB Connected (tenda)")
	return db
}