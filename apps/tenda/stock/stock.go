package stock

import (
	"fmt"
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

type Picked struct {
	PID          int 
	PNO          int
	Model        string
	Qty          int
	Customer     string
	Last_updated string
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

func GetModel(querymodel string) []Model {
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
        if err != nil {
			fmt.Println("DB Query failed error: ", err)
		}
        allmodels = append(allmodels, m)
    }
	return allmodels
}
func GetLocationModels(querylocation string) []Model {
	sql := "SELECT location, model, unit, cartons, boxes, total " +
		   "FROM stock_locations " +
		   "WHERE location='" + querylocation + "'"
	rows, err := db.Query(sql)
	fmt.Println("SQL Statement:", sql)
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	m := Model{}
	allmodels := []Model{}
    for rows.Next() {
        err := rows.Scan(&m.Location, &m.Model, &m.Unit, &m.Cartons, &m.Boxes, &m.Total)
        if err != nil {
			fmt.Println("DB Query failed error: ", err)
		}
        allmodels = append(allmodels, m)
    }
    fmt.Printf("allmodels:%v\n", allmodels)
	return allmodels
}

func InsertPicked(PNO string, model string, qty string, customer string, updated string) int64  {
	
	stmt, err := db.Prepare("INSERT INTO picked(PNO, model, qty, customer, last_updated) VALUES (?,?,?,?,?)")
	if err != nil {
		fmt.Println("Error sql Prepare:", err)
	}
	res, err := stmt.Exec(PNO, model, qty, customer, updated)
	if err != nil {
		fmt.Println("Error exectue sql:", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		fmt.Println("Error last ID:", err)
	}

	fmt.Println("Last id:", id)
	return id
}

func GetTodayPackages(date string) []Picked {
	sql := "SELECT PID, PNO, model, qty, customer, last_updated " +
		   "FROM picked " +
		   "WHERE last_updated LIKE '%" + date + "%'"
	rows, err := db.Query(sql)
	fmt.Println("SQL Statement:", sql)
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	p := Picked{}
	allPicked := []Picked{}
    for rows.Next() {
        err := rows.Scan(&p.PID, &p.PNO, &p.Model, &p.Qty, &p.Customer, &p.Last_updated)
        if err != nil {
			fmt.Println("DB Query failed error: ", err)
		}
        allPicked = append(allPicked, p)
    }
    fmt.Printf("allPicked:%v\n", allPicked)
	return allPicked
}
