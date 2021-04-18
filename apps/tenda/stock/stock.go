package stock

import (
	"fmt"
	"github.com/guangxue/webpages/mysql"
	_ "github.com/go-sql-driver/mysql"	
)
var db = mysql.Connect("tenda");

type model struct {
	Model    string `json:"model"`
	Location string `json:"location"`
	Unit     int    `json:"unit"`
	Cartons  int    `json:"cartons"`
	Boxes    int    `json:"loose"`
	Total    int    `json:"total"`
}

func GetModelLocation(modelName string) {

}

func GetAllModels() []string {
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
		models = append(models, thisModel)
	}
	return models
}

func GetModel(querymodel string) []model {
	sql := "SELECT model, location, unit, cartons, boxes, total " +
               "FROM stock_locations " +
               "WHERE model ='" + querymodel + "'"
    rows, err := db.Query(sql)
    if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	m := model{}
    allmodels := []model{}
    for rows.Next() {
        err := rows.Scan(&m.Model, &m.Location, &m.Unit, &m.Cartons, &m.Boxes, &m.Total)
        if err !=nil {

        }
        allmodels = append(allmodels, m)
    }
	return allmodels
}