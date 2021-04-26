package mysql

import (
	"fmt"
	"database/sql"
	"strings"
	"strconv"
	_ "github.com/go-sql-driver/mysql"	
)

type Statement struct {
	SelectColumns string
	ColumnCount   int
	ColumnNames   []string
	TableName     string
	InsertStmt    string
	InsertValues  []interface{}
	SetExpr  	  string
	WhereClause   string
	DBconnection  *sql.DB
	QueryType  string
}

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
	fmt.Println("[DB] Connected")
	return db
}

func Select(searchColumns ...string) *Statement{
	sqlstmt := &Statement{}
	sqlstmt.ColumnCount = len(searchColumns)
	sqlstmt.ColumnNames = searchColumns
	sqlstmt.QueryType = "SELECT"
	sqlstmt.SelectColumns = "SELECT " + strings.Join(searchColumns, ", ")
	return sqlstmt
}
func SelectDistinct(searchColumns ...string) *Statement{
	sqlstmt := &Statement{}
	sqlstmt.ColumnCount = len(searchColumns)
	sqlstmt.ColumnNames = searchColumns
	sqlstmt.QueryType = "SELECT"
	sqlstmt.SelectColumns = "SELECT DISTINCT " + strings.Join(searchColumns, ", ")
	return sqlstmt
}

func InsertInto(tableName string, insertColumns []string, insertValues []interface{}) *Statement{
	sqlstmt := &Statement{}
	sqlstmt.QueryType = "INSERT"
	sqlstmt.InsertValues = insertValues

	placeholders := make([]string, len(insertValues))
	for i, _ := range placeholders {
		placeholders[i] = "?"
	}
	insertPlaceholders := "("+strings.Join(placeholders, ",")+")"
	fmt.Println("INSERT placeholders:", insertPlaceholders)
	insertStmt := "INSERT INTO "+tableName + "("+strings.Join(insertColumns, ",")+") VALUES " + insertPlaceholders
	sqlstmt.InsertStmt = insertStmt
	fmt.Println("insertStmt:", insertStmt)
	return sqlstmt
}

func Update(tableName string) *Statement{
	sqlstmt := &Statement{}
	sqlstmt.TableName = "Update "+tableName
	sqlstmt.QueryType = "UPDATE"
	return sqlstmt
}

func (sqlstmt *Statement) Set(updateColumns map[string]string) *Statement{
	setExpression := " SET "
	for col, val := range updateColumns {
		setExpression += col + "='" +val+ "', "
	}
	sqlstmt.SetExpr = setExpression[:len(setExpression)-2]
	return sqlstmt
}

func (sqlstmt *Statement) From(tableName string) *Statement{
	sqlstmt.TableName = " FROM "+tableName
	return sqlstmt
}

func (sqlstmt *Statement) Where(column string, searchColumn string) *Statement{
	sqlstmt.WhereClause = " Where " + column + "='" + searchColumn + "'"
	return sqlstmt
}

func (sqlstmt *Statement) WhereBetween(value1 string, value2 string) *Statement{
	sqlstmt.WhereClause = " Where BETWEEN" + value1 + " AND " + value2
	return sqlstmt
}

func (sqlstmt *Statement) WhereLike(column string, pattern string) *Statement{
	sqlstmt.WhereClause = " Where " + column + " LIKE '" + pattern + "'"
	return sqlstmt
}

func (sqlstmt *Statement) Use(db *sql.DB) []map[string]string{
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
	fmt.Println("QueryType:", sqlstmt.QueryType)
	switch sqlstmt.QueryType {
	case "SELECT":
		
		// columnsToSelect := strings.Join(searchColumns, ", ")
		stmt := sqlstmt.SelectColumns + sqlstmt.TableName + sqlstmt.WhereClause
		fmt.Printf("final stmt: %s\n", stmt)
		var scannedColumns = make([]interface{}, sqlstmt.ColumnCount)
		
		// convert []interface{} to slice -> for easing indexing with [1]
		// save each interface{} with string poiner -> for rows.Scan()
		for idx, _ := range sqlstmt.ColumnNames {
			scannedColumns[idx] = new(string)
		}
		rows, err := db.Query(stmt)
		if err != nil {
			fmt.Println("[stmt.Get] db.Query error line:135:", err)
		}
		for rows.Next() {
			err := rows.Scan(scannedColumns...)
			if err != nil {
				fmt.Println("[Get]: dbColumns Scan error:140:", err)
			}
			// save each scanned column to col map[string]string
			col := map[string]string{}
			for idx, val := range sqlstmt.ColumnNames {
				colstr, ok := scannedColumns[idx].(*string)
				if !ok {
					fmt.Println("Cannot convert *interface{} to *string")
				}
				col[val] = *colstr
			}
			// append scanned column{map} to slice of maps
			finalColumns = append(finalColumns, col)
		}
		
	case "UPDATE":
		stmt := sqlstmt.TableName + sqlstmt.SetExpr + sqlstmt.WhereClause
		fmt.Println("UPdate Statement:", stmt)
		

	case "INSERT":
		fmt.Println("INSERT statment::", sqlstmt.InsertStmt)
		fmt.Println("INSERT values:", sqlstmt.InsertValues)
		stmt, err := db.Prepare(sqlstmt.InsertStmt)
		if err != nil {
			fmt.Println("Error sql Prepare:", err)
		}
		res, err := stmt.Exec(sqlstmt.InsertValues...)
		if err != nil {
			fmt.Println("Error exectue sql:", err)
		}

		id, err := res.LastInsertId()

		if err != nil {
			fmt.Println("Error last ID:", err)
		}

		fmt.Println("[InsertPicked] Last id:", id)
		rid := strconv.FormatInt(id, 10)
		insertFeedback := map[string]string{"lastId":rid}
		finalColumns = append(finalColumns, insertFeedback)
		

	} // EOS: end of switch

	for _,val := range finalColumns {
		fmt.Println("stmt: finalColumns:", val)
	}
	return finalColumns
}


