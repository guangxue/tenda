package mysql

import (
	"context"
	"fmt"
	"database/sql"
	"strings"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
)

type Statement struct {
	SelectColumns  string
	ColumnCount    int
	ColumnNames    []string
	TableName      string
	InsertStmt     string
	InsertValues   []interface{}
	SetExpr	       string
	WhereClause    string
	AndWhereClause string
	QueryType      string
	UpdateNoWhere  bool
	DeleteNoWhere  bool
	RawStatment    string
}

func Connect(dbname string) *sql.DB {
	dsn := "gzhang:guangxue@/" + dbname
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	var pingerr = db.Ping()
	if pingerr != nil {
		fmt.Println("DB pinging error: ", pingerr)
	}

    // fmt.Printf("[%-18s] Running on `dev` branch\n", "...")
	fmt.Printf("[%-18s] Connected\n", " -- mysql.go")
	return db
}

func Begin(db *sql.DB) (*sql.Tx, context.Context) {
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		fmt.Println("DB BeginTx error: ", err)
	}
	return tx, ctx
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

func InsertInto(tableName string, insertQuery map[string]interface{}) *Statement{
	sqlstmt := &Statement{}
	sqlstmt.QueryType = "INSERT"

	placeholders := make([]string, len(insertQuery))
	for i, _ := range placeholders {
		placeholders[i] = "?"
	}
	insertPlaceholders := "("+strings.Join(placeholders, ",")+")"
	// fmt.Println("INSERT placeholders:", insertPlaceholders)
	insertColumns := []string{}
	insertValues  := []interface{}{}
	for col, val := range insertQuery {
		insertColumns = append(insertColumns, col)
		insertValues  = append(insertValues, val)
	}
	sqlstmt.InsertValues = insertValues
	insertStmt := "INSERT INTO "+tableName + "("+strings.Join(insertColumns, ",")+") VALUES " + insertPlaceholders
	sqlstmt.InsertStmt = insertStmt
	// fmt.Println("insertStmt:", insertStmt)
	return sqlstmt
}

func Update(tableName string, updateNoWhere bool) *Statement{
	sqlstmt := &Statement{}
	sqlstmt.UpdateNoWhere = updateNoWhere
	sqlstmt.TableName = "UPDATE "+tableName
	sqlstmt.QueryType = "UPDATE"
	return sqlstmt
}

func DeleteFrom(tableName string, deleteNoWhere bool) *Statement {
	sqlstmt := &Statement{}
	sqlstmt.DeleteNoWhere = deleteNoWhere
	sqlstmt.TableName = "DELETE FROM "+tableName
	sqlstmt.QueryType = "DELETE"
	return sqlstmt
}

func SelectRaw(rawStmt string, columnNames ...string) *Statement {
	sqlstmt := &Statement{}
	sqlstmt.RawStatment = rawStmt
	sqlstmt.ColumnNames = columnNames
	sqlstmt.ColumnCount = len(columnNames)
	sqlstmt.QueryType = "SELECTRAW"
	return sqlstmt
}

func (sqlstmt *Statement) Set(updateColumns map[string]interface{}) *Statement{
	setExpression := " SET "
	for col, val := range updateColumns {
		setExpression += col + "='" + fmt.Sprintf("%v",val) + "', "
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

func (sqlstmt *Statement) WhereBetween(col string,value1 string, value2 string) *Statement{
	sqlstmt.WhereClause = " Where "+ col +" BETWEEN " + value1 + " AND " + value2
	return sqlstmt
}

func (sqlstmt *Statement) WhereLike(column string, pattern string) *Statement{
	sqlstmt.WhereClause = " Where " + column + " LIKE '" + pattern + "'"
	return sqlstmt
}

func (sqlstmt *Statement) AndWhere(column string, condition string, searchColumn string) *Statement {
	if len(sqlstmt.WhereClause) > 0 {
		sqlstmt.AndWhereClause += " AND " + column + " " +condition + " '"+searchColumn + "'"
		return sqlstmt
	}
	return sqlstmt
}

func (sqlstmt *Statement) Use(db *sql.DB) []map[string]string{
	
	finalColumns := []map[string]string{}
	// fmt.Printf("[QueryType] %s\n", sqlstmt.QueryType)
	switch sqlstmt.QueryType {
	case "SELECT", "SELECTRAW":
		
		// columnsToSelect := strings.Join(searchColumns, ", ")
		stmt := ""
		if sqlstmt.QueryType == "SELECT" {
			stmt = sqlstmt.SelectColumns + sqlstmt.TableName + sqlstmt.WhereClause + sqlstmt.AndWhereClause
			fmt.Printf("[%-18s] %s\n", " -- SELECT --", sqlstmt.SelectColumns)
			fmt.Printf("[%-18s]  %s\n", "", sqlstmt.TableName)
			fmt.Printf("[%-18s]  %s\n", "", sqlstmt.WhereClause)
			fmt.Printf("[%-18s]  %s\n", " -- END ----", sqlstmt.AndWhereClause)

		} else {
			stmt = sqlstmt.RawStatment
			fmt.Printf("[%-18s] %s\n", "SelectRaw", sqlstmt.RawStatment)
		}
		
		var scannedColumns = make([]interface{}, sqlstmt.ColumnCount)
		
		/* convert []interface{} to slice -> for easing indexing with [1]
		 | save each interface{} with string poiner -> for rows.Scan()
		 */
		for idx, _ := range sqlstmt.ColumnNames {
			scannedColumns[idx] = new(string)
		}
		rows, err := db.Query(stmt)
		if err != nil {
			fmt.Println("[db *Err] tx.QueryContext error:", err)
		}
		for rows.Next() {
			err := rows.Scan(scannedColumns...)
			if err != nil {
				fmt.Println("[db *Err]: rows.Scan error:", err)
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
		if sqlstmt.UpdateNoWhere {
			stmt := sqlstmt.TableName + sqlstmt.SetExpr
			fmt.Printf("[%-18s] %s\n", "-- UPDATE --", sqlstmt.TableName)
			fmt.Printf("[%-18s]  %s\n", " SET", sqlstmt.SetExpr)
			res, err := db.Exec(stmt)
			if err != nil {
				fmt.Println("[db *Err]: Update error:", err)
			}
			rnums, err := res.RowsAffected()
			if err != nil {
				fmt.Println("[db *Err] RowsAffected:", err)
			}
			fmt.Printf("[%-18s] %d\n", "UPDATE rows", rnums)
			rid := strconv.FormatInt(rnums, 10)
			rowsFeedback := map[string]string{"rowsAffected":rid}
			finalColumns = append(finalColumns, rowsFeedback)
		} else if len(sqlstmt.WhereClause) > 0 {
			stmt := sqlstmt.TableName + sqlstmt.SetExpr + sqlstmt.WhereClause + sqlstmt.AndWhereClause
			fmt.Printf("[%-18s] %s\n",  " -- UPDATE --", sqlstmt.TableName)
			fmt.Printf("[%-18s]  %s\n", ".. SET", sqlstmt.SetExpr)
			fmt.Printf("[%-18s]  %s\n", ".. WHERE", sqlstmt.WhereClause)
			fmt.Printf("[%-18s]  %s\n", ".. AND", sqlstmt.AndWhereClause)
			
			res, err := db.Exec(stmt)
			if err != nil {
				fmt.Println("[db *Err]: Update error:", err)
			}
			rnums, err := res.RowsAffected()
			if err != nil {
				fmt.Println("[db *Err] RowsAffected:", err)
			}
			fmt.Printf("[%-18s] %d\n", "UPDATE rows", rnums)
			rid := strconv.FormatInt(rnums, 10)
			rowsFeedback := map[string]string{"rowsAffected":rid}
			finalColumns = append(finalColumns, rowsFeedback)
		} else {
			// fmt.Printf(">> %s\n", stmt)
			fmt.Println("[db *Err] WhereClause needed!")
		}
	case "INSERT":
		fmt.Printf("[%-18s] %s\n", "INSERT",sqlstmt.InsertStmt)
		fmt.Printf("[%-18s] %v\n", "INSERT VALUES", sqlstmt.InsertValues)
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

		fmt.Printf("[%-18s] %s: %d\n", "INSERT","Last Insert Id",id)
		rid := strconv.FormatInt(id, 10)
		insertFeedback := map[string]string{"lastId":rid}
		finalColumns = append(finalColumns, insertFeedback)
	case "DELETE":
		if sqlstmt.DeleteNoWhere {
			fmt.Printf("[%-18s] %s\n", "DELETE","!!! Deleting with NO WHERE !!!")
		} else if len(sqlstmt.WhereClause) > 0 {
			stmt := sqlstmt.TableName + sqlstmt.SetExpr + sqlstmt.WhereClause + sqlstmt.AndWhereClause
			fmt.Printf("[%-18s] %s: %s\n", "DELETE","stmt:",stmt)
			res, err := db.Exec(stmt)
			if err != nil {
				fmt.Println("[db *Err]: Delete error:", err)
			}
			rnums, err := res.RowsAffected()
			if err != nil {
				fmt.Println("[db *Err] RowsAffected:", err)
			}
			fmt.Printf("[%-18s] %d\n", "Delete rows", rnums)
			rid := strconv.FormatInt(rnums, 10)
			rowsFeedback := map[string]string{"rowsAffected":rid}
			finalColumns = append(finalColumns, rowsFeedback)
		} else {
			fmt.Printf("[%-18s] %s: %d\n", "DELETE","WRONG SQL Stmt")
		}
	} // EOS: end of switch

	return finalColumns
}

func (sqlstmt *Statement) With(tx *sql.Tx, ctx context.Context) []map[string]string{
	
	finalColumns := []map[string]string{}
	// fmt.Printf("[QueryType] %s\n", sqlstmt.QueryType)
	switch sqlstmt.QueryType {
	case "SELECT", "SELECTRAW":
		
		// columnsToSelect := strings.Join(searchColumns, ", ")
		stmt := ""
		if sqlstmt.QueryType == "SELECT" {
			stmt = sqlstmt.SelectColumns + sqlstmt.TableName + sqlstmt.WhereClause + sqlstmt.AndWhereClause
			fmt.Printf("[%-18s] %s\n", "SELECT", sqlstmt.SelectColumns)
			fmt.Printf("[%-18s]  %s\n", "SELECT FROM", sqlstmt.TableName)
			fmt.Printf("[%-18s]  %s\n", "SELECT WHERE", sqlstmt.WhereClause)
			fmt.Printf("[%-18s]  %s\n", "SELECT AND", sqlstmt.AndWhereClause)

		} else {
			stmt = sqlstmt.RawStatment
			fmt.Printf("[%-18s] %s\n", "SelectRaw", sqlstmt.RawStatment)
		}
		
		var scannedColumns = make([]interface{}, sqlstmt.ColumnCount)
		
		/* convert []interface{} to slice -> for easing indexing with [1]
		 | save each interface{} with string poiner -> for rows.Scan()
		 */
		for idx, _ := range sqlstmt.ColumnNames {
			scannedColumns[idx] = new(string)
		}
		rows, err := tx.QueryContext(ctx, stmt)
		if err != nil {
			tx.Rollback()
			fmt.Println("[db *Err] tx.QueryContext error:", err)
		}
		for rows.Next() {
			err := rows.Scan(scannedColumns...)
			if err != nil {
				fmt.Println("[db *Err]: rows.Scan error:", err)
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
		if sqlstmt.UpdateNoWhere {
			stmt := sqlstmt.TableName + sqlstmt.SetExpr
			fmt.Printf("[%-18s] %s\n", "UPDATE", sqlstmt.TableName)
			fmt.Printf("[%-18s]  %s\n", "UPDATE SET", sqlstmt.SetExpr)
			res, err := tx.ExecContext(ctx, stmt)
			if err != nil {
				tx.Rollback()
				fmt.Println("[db *Err]: Update error:", err)
			}
			rnums, err := res.RowsAffected()
			if err != nil {
				fmt.Println("[db *Err] RowsAffected:", err)
			}
			fmt.Printf("[%-18s] %d\n", "UPDATE rows", rnums)
			rid := strconv.FormatInt(rnums, 10)
			rowsFeedback := map[string]string{"rowsAffected":rid}
			finalColumns = append(finalColumns, rowsFeedback)
		} else if len(sqlstmt.WhereClause) > 0 {
			stmt := sqlstmt.TableName + sqlstmt.SetExpr + sqlstmt.WhereClause + sqlstmt.AndWhereClause
			fmt.Printf("[%-18s] \n",  "* BEGIN Transaction *")
			fmt.Printf("[%-18s] %s\n",  "UPDATE", sqlstmt.TableName)
			fmt.Printf("[%-18s]  %s\n", "UPDATE SET", sqlstmt.SetExpr)
			fmt.Printf("[%-18s]  %s\n", "UPDATE WHERE", sqlstmt.WhereClause)
			fmt.Printf("[%-18s]  %s\n", "UPDATE AND", sqlstmt.AndWhereClause)
			
			res, err := tx.ExecContext(ctx, stmt)
			if err != nil {
				tx.Rollback()
				fmt.Println("[db *Err]: Update error:", err)
			}
			rnums, err := res.RowsAffected()
			if err != nil {
				fmt.Println("[db *Err] RowsAffected:", err)
			}
			fmt.Printf("[%-18s] %d\n", "UPDATE rows", rnums)
			rid := strconv.FormatInt(rnums, 10)
			rowsFeedback := map[string]string{"rowsAffected":rid}
			finalColumns = append(finalColumns, rowsFeedback)
		} else {
			// fmt.Printf(">> %s\n", stmt)
			fmt.Println("[db *Err] WhereClause needed!")
		}
	case "INSERT":
		fmt.Printf("[%-18s] \n",  "* BEGIN Transaction *")
		fmt.Printf("[%-18s] %s\n", "INSERT",sqlstmt.InsertStmt)
		fmt.Printf("[%-18s] %v\n", "INSERT VALUES", sqlstmt.InsertValues)
		stmt, err := tx.PrepareContext(ctx, sqlstmt.InsertStmt)
		// stmt, err := db.Prepare(sqlstmt.InsertStmt)
		if err != nil {
			fmt.Println("Error sql Prepare:", err)
		}
		res, err := tx.StmtContext(ctx, stmt).Exec(sqlstmt.InsertValues...)
		if err != nil {
			tx.Rollback()
			fmt.Println("Error exectue sql:", err)
		}

		id, err := res.LastInsertId()

		if err != nil {
			fmt.Println("Error last ID:", err)
		}

		fmt.Printf("[%-18s] %s: %d\n", "INSERT","Last Insert Id",id)
		rid := strconv.FormatInt(id, 10)
		insertFeedback := map[string]string{"lastId":rid}
		finalColumns = append(finalColumns, insertFeedback)
	case "DELETE":
		if sqlstmt.DeleteNoWhere {
			fmt.Printf("[%-18s] %s\n", "DELETE","!!! Deleting with NO WHERE !!!")
		} else if len(sqlstmt.WhereClause) > 0 {
			
			stmt := sqlstmt.TableName + sqlstmt.SetExpr + sqlstmt.WhereClause + sqlstmt.AndWhereClause
			fmt.Printf("[%-18s] %s: %s\n", "DELETE","stmt:",stmt)
			res, err := tx.ExecContext(ctx, stmt)
			// res, err := db.Exec(stmt)
			if err != nil {
				fmt.Println("[db *Err]: Delete error:", err)
				tx.Rollback()
			}
			rnums, err := res.RowsAffected()
			if err != nil {
				fmt.Println("[db *Err] RowsAffected:", err)
			}
			fmt.Printf("[%-18s] %d\n", "Delete rows", rnums)
			rid := strconv.FormatInt(rnums, 10)
			rowsFeedback := map[string]string{"rowsAffected":rid}
			finalColumns = append(finalColumns, rowsFeedback)
		} else {
			fmt.Printf("[%-18s] %s: %d\n", "DELETE","WRONG SQL Stmt")
		}
	} // EOS: end of switch

	return finalColumns
}
