package models

import (
	"fmt"
	"database/sql"
	"strings"
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

type WhereClause struct {
	ID       string
	Model    string
	Location string
	Last_update string
}

func searchConditions(search WhereClause) string {
	searchCondition := ""
	switch {
	case len(search.ID) > 0:
		searchCondition = " WHERE SID='" + search.ID +  "'"
		return searchCondition
	case len(search.Model) > 0:
		searchCondition = " WHERE model='" + search.Model +  "'"
		return searchCondition
	case len(search.Location) > 0:
		searchCondition = " WHERE location='" + search.Location +  "'"
		return searchCondition
	default:
		return ""
	}
}

func Names() []string {
	rows, err := db.Query("select distinct model from stock_locations")
	if err != nil {
		fmt.Println("DB Query failed error: ", err)
	}
	name := ""
	modelNames := []string{}
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {

		}
		modelNames = append(modelNames, name)
	}
	return modelNames
}
func GetAll(search WhereClause) []Model {
	sql :=  "SELECT model, location, unit, cartons, boxes, total FROM stock_locations" + searchConditions(search)
    rows, err := db.Query(sql)
    if err != nil {
		fmt.Println("[getAllModel] DB Query error: ", err)
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


func Get(tableName string, search WhereClause, searchColumns ...string) []map[string]string{
/***
	returned finalColumns should be slices of maps:
	[
		{model:MW6-3PK, location: 1-G-1},
		{model:MW6-3PK, location: 5-G-1},
		{model:MW6-3PK, location: 7-3-1},
		{model:MW6-3PK, location: 8-1-2},
		{model:MW6-3PK, location: 0-G-5},
	]
****/
	finalColumns := []map[string]string{}
	columnsToSelect := strings.Join(searchColumns, ", ")
	sqlstmt := "SELECT " + columnsToSelect + " FROM " + tableName + searchConditions(search)
	fmt.Println("[Get sqlstmt]: ", sqlstmt)
	var scannedColumns = make([]interface{}, len(searchColumns))
	
	// convert []interface{} to slice -> for easing indexing with [1]
	// save each interface{} with string poiner -> for rows.Scan()
	for idx, _ := range searchColumns {
		scannedColumns[idx] = new(string)
	}
	rows, err := db.Query(sqlstmt)
	if err != nil {
		fmt.Println("[Get] Query error:104:", err)
	}
	for rows.Next() {
		err := rows.Scan(scannedColumns...)
		if err != nil {
			fmt.Println("[Get]: dbColumns Scan error:109:", err)
		}
		// save each scanned column to col map[string]string
		col := map[string]string{}
		for idx, val := range searchColumns {
			colstr, ok := scannedColumns[idx].(*string)
			if !ok {
				fmt.Println("Cannot convert *interface{} to *string")
			}
			col[val] = *colstr
		}
		// append scanned column{map} to slice of maps
		finalColumns = append(finalColumns, col)
	}
	return finalColumns
}

/*
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

*/