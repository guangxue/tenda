package stock

import (
	"fmt"
	"database/sql"
	"github.com/guangxue/webapps/mysql"
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
	PID          int            `json:"PID"`
	PNO          string         `json:"PNO"`
	Model        string         `json:"model"`
	Qty          int            `json:"qty"`
	Customer     sql.NullString `json:"customer"`
	Location     string         `json:"location"`
	Status       string         `json:"status"`
	Last_updated string         `json:"updated"`
}

type Loc struct {
	Location string `json:"location"`
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

func InsertPicked(PNO string, model string, qty string, customer string, location string, status string, updated string) int64  {
	stmt, err := db.Prepare("INSERT INTO picked(PNO, model, qty, customer, location, status, last_updated) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("Error sql Prepare:", err)
	}
	res, err := stmt.Exec(PNO, model, qty, customer, location, status, updated)
	if err != nil {
		fmt.Println("Error exectue sql:", err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		fmt.Println("Error last ID:", err)
	}

	fmt.Println("[InsertPicked] Last id:", id)
	return id
}

func GetTodayPackages(date string) []Picked {
	sql := "SELECT PID, PNO, model, qty, customer, location, status, last_updated " +
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
        err := rows.Scan(&p.PID, &p.PNO, &p.Model, &p.Qty, &p.Customer, &p.Location, &p.Status, &p.Last_updated)
        if err != nil {
			fmt.Println("DB Query failed error: ", err)
		}
        allPicked = append(allPicked, p)
    }
    fmt.Printf("[DB selected] allPicked:%v\n", len(allPicked))
	return allPicked
}


func GetModelLocations(model string) []Loc {
	sql := "SELECT location from stock_locations where model='"+model+"'"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	l := Loc{}
	locs := []Loc{}
	for rows.Next() {
		err := rows.Scan(&l.Location)
		if err != nil {
			fmt.Println("DB location query failed:", err)
		}
		locs = append(locs, l)
	}
	fmt.Printf("[GetModelLocations] locations array: %v\n", locs)
	return locs

}
func GetPendingParcels() []Picked {
	sql := "SELECT PID, PNO, model, qty, customer, location, status, last_updated " +
		   "FROM picked " +
		   "WHERE status='Pending'"
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	p := Picked{}
	allPicked := []Picked{}
    for rows.Next() {
        err := rows.Scan(&p.PID, &p.PNO, &p.Model, &p.Qty, &p.Customer, &p.Location, &p.Status, &p.Last_updated)
        if err != nil {
			fmt.Println("DB Query failed error: ", err)
		}
        allPicked = append(allPicked, p)
    }
    fmt.Printf("rows.Scaned - allPicked:%v\n", allPicked)
	return allPicked
}

