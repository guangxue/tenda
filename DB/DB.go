package DB

import  (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Statement struct {
	SelectQuery string
	WhereQuery  string
	UpdateQuery string
	TableName   string
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
	fmt.Println("DB Connected (tenda)")
	return db
}

func Table(tableName string) *Statement{
	var stmt = &Statement{}
	stmt.TableName = tableName
	return stmt
}

func(s *Statement) Select(columns string) *Statement {
	s.SelectQuery = "SELECT " + columns
	return s;
}

func (s *Statement) Get() {
	fmt.Println("select statement:", s.SelectQuery)
}

func(s *Statement) Update(tableName string) *Statement {
	s.UpdateQuery = "Update " + tableName
	return s;
}
