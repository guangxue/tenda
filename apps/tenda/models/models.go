package models

import (
	"fmt"
	"database/sql"
	"github.com/guangxue/webpages/mysql"
	_ "github.com/go-sql-driver/mysql"	
)
var db = mysql.Connect("tenda");

type Model struct {
	Model    string `json:"model"`
	Location string `json:"location"`
	Unit     int    `json:"unit"`
	Cartons  int    `json:"cartons"`
	Boxes    int    `json:"loose"`
	Total    int    `json:"total"`
}


func GetAll() {
	rows, err := db.Query("select distinct model from stock_locations")
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	thisModel := ""
	models := []string{}
	for rows.Next() {
		err := rows.Scan(&thisModel)
		if err != nil {

		}
		models = append(models, modelName)
	}
	return models
}

func GetQueryModel(querymodel string) {
	sql := "SELECT model, location, unit, cartons, boxes, total " +
               "FROM stock_locations " +
               "WHERE model ='" + querymodel + "'"
    rows, err := db.Query(sql)
    if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	m := Model{}
    allmodels := []Model{}
    for rows.Next() {
        err := rows.Scan(&m.Model, &m.Location, &m.Unit, &m.Cartons, &m.Boxes, &m.Total)
        ErrorCheck(err)
        allmodels = append(allmodels, m)
    }
	return allmodels
}