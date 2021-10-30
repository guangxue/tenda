package tenda

import (
	"database/sql"
	"fmt"
	"time"
	"html/template"
	"net/http"
	"encoding/json"
	"github.com/guangxue/webapps/mysql"
)

var tbname map[string]string = map[string]string{}
var dbCommits map[string]*sql.Tx = map[string]*sql.Tx{}

type pickcolumns struct {
	PID      int
	Model    string
	Qty      int
	Location string
}
type updateModel struct {
	Location string
	Model    string
	Unit     int
	Cartons  int
	Boxes    int
	Total    int
}

var db = mysql.Connect("tenda");
func init() {
	tbname["stock_updated"] = "stock_updated"
	tbname["last_updated"] = "last_updated"
	tbname["picklist"] = "picklist"

	// tbname["stock_updated"] = "stock_updated_test"
	// tbname["last_updated"] = "last_updated_test"
	// tbname["picklist"] = "picklist_test"
}

func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func returnJson(w http.ResponseWriter, data []map[string]string) {
	jsn, err := json.Marshal(data)
	if err != nil {
		fmt.Println("data JSON Marshal error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}

func returnJs(w http.ResponseWriter, data map[string]string) {
	jsn, err := json.Marshal(data)
	if err != nil {
		fmt.Println("data JSON Marshal error: ", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}

func TimeNow() string {
	ts := time.Now()
	now := ts.Format("2006-01-02 15:04:05")
	return now
}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	tmplpath := "templates/tenda/" + templateName
	tmpl, err := template.ParseFiles(tmplpath, "templates/tenda/base.html","templates/tenda/nav.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}

func RenderHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmplpath := "templates/tenda/" + templateName
		tmpl, err := template.ParseFiles("templates/tenda/base.html", "templates/tenda/nav.html", tmplpath)
		if err != nil {
			fmt.Println("template parsing errors: ", err)
		}
		err = tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			fmt.Println("template executing errors: ", err)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	render(w, "login.html", nil)
}

func Locations(w http.ResponseWriter, r *http.Request) {
	searchModel := r.URL.Query().Get("model");
	fmt.Printf("[ **%-12s**:] %s\n", " Fetch URL", "/tenda/api/locations")
	fmt.Printf("[%-18s] ?model = %s\n", " -- tenda.go", searchModel)
	if len(searchModel) == 0 {
		emptyLocation := make(map[string]string)
		returnJs(w, emptyLocation)
	}
	switch searchModel {
	case "MW6-3PK":
		fmt.Println("====== searchModel: ==", searchModel)
		locs := []map[string]string{}
		loc := map[string]string{
			"location":"0-G-4",
		}
		locs = append(locs, loc)
		returnJson(w, locs)
	default:
		allLocations := mysql.
			Select("location").
			From(tbname["stock_updated"]).
			Where("model", searchModel).
			AndWhere("total", ">", "0").
		Use(db)
		for _, locs := range allLocations {
			for key, val := range locs {
				fmt.Println("key=%s, val=%s\n", key, val)
			}
		}
		returnJson(w, allLocations)
	}
	// if len(searchModel) > 0 {
	// 	allLocations := mysql.
	// 		Select("location").
	// 		From(tbname["stock_updated"]).
	// 		Where("model", searchModel).
	// 		AndWhere("total", ">", "0").
	// 	Use(db)
	// 	returnJson(w, allLocations)
	// }
}

func PickListUpdatePage(w http.ResponseWriter, r *http.Request) {

	dbPickedInfo := map[string]string{}

	if r.Method == http.MethodGet {
		queryPID := r.URL.Query().Get("PID")
		LID := r.URL.Query().Get("LID")

		if queryPID != "" {
			fmt.Printf("[%-18s] PID:%s\n", "PickListUpdatePage",queryPID)
			currentPID := mysql.Select("PNO", "model", "sales_mgr","qty", "customer", "location", "status").From(tbname["picklist"]).Where("PID", queryPID).Use(db)
			dbPickedInfo = currentPID[0]
			dbPickedInfo["PID"] = queryPID
			dbPickedInfo["status2"] = ""
			if dbPickedInfo["status"] == "Pending" {
				dbPickedInfo["status2"] = "Complete"
			} else {
				dbPickedInfo["status2"] = "Pending"
			}
			fmt.Printf("[%-18s] /picklist GET PID:%s\n", "UpdatePickList", queryPID)
			fmt.Printf("[%-18s] return:%s\n", "UpdatePickList", dbPickedInfo)
			render(w, "picklistupdate.html", dbPickedInfo)
		}
		if LID != "" {
			currentLID := mysql.
				Select("LID", "location", "model", "unit", "cartons", "boxes", "total", "completed_at").
				From(tbname["last_updated"]).
				Where("LID", LID).
			Use(db)
			returnJson(w, currentLID)
		}
	}
}

func PickListInspectPage(w http.ResponseWriter, r *http.Request) {
	modelName := r.URL.Query().Get("model")
	location  := r.URL.Query().Get("location")
	pickDate  := r.URL.Query().Get("pickDate")
	if modelName == "" || location == "" || pickDate == "" {
		infoJs := map[string]string {
			"error": "empty modelName,location,pickDate",
		}
		returnJs(w, infoJs)
	}
	fmt.Println("model name", modelName)
	fmt.Println("location", location)
	fmt.Println("pickDate", pickDate)
	startDate := fmt.Sprintf("'%s'", pickDate)
	betweenDate := fmt.Sprintf("date_add('%s', interval 7 day)", pickDate)
	inspectModel := mysql.
		Select("PID","PNO", "sales_mgr", "model", "qty", "customer", "location", "created_at").
		From(tbname["picklist"]).
		WhereBetween("created_at", startDate, betweenDate).
		AndWhere("model", "=", modelName).
		AndWhere("location", "=", location).
		AndWhere("status", "=", "Complete").
		Use(db)

	fmt.Println("Inspected Model")
	for i, model := range inspectModel {
		fmt.Printf("------------[%d]---------\n", i)
		for key, val := range model {
			fmt.Printf("\t%s : %s\n", key, val)
		}
	}
	render(w, "inspectpicklist.html", inspectModel)
}

func StockUpdatePage(w http.ResponseWriter, r *http.Request) {
	SID := r.URL.Query().Get("SID")

	if SID != "" {
		currentStockToUpdate := mysql.
			Select("SID", "location", "model", "unit","kind", "cartons", "boxes","total", "update_comments").
			From(tbname["stock_updated"]).
			Where("SID", SID).
			Use(db)
		fmt.Println("currentStockToUpdate:", currentStockToUpdate)
		render(w, "stockupdate.html", currentStockToUpdate[0])
	}
}

func LastUpdatedPage(w http.ResponseWriter, r *http.Request) {
	LID := r.URL.Query().Get("LID");
	fmt.Println("[LastUpdatedPage] LID:", LID)
	if LID != "" {
		LastUpdated := mysql.
			Select("LID","location","model","unit","old_total","total_picks","cartons","boxes","completed_at").
			From(tbname["last_updated"]).
			Where("LID", LID).
			Use(db)
		fmt.Println("last_updated:", LastUpdated[0])
		render(w, "lastupdated.html", LastUpdated[0])
	}
}


func TxCommit(w http.ResponseWriter, r *http.Request) {
	commitName := r.URL.Query().Get("cmn")
	resText    := map[string]string{}
	
	fmt.Printf("[%-18s] Commit name : %s\n", "TxCommit",commitName)
	fmt.Println("[* END Transaction *]")

	tx, ok := dbCommits[commitName]
	if !ok {
		fmt.Println("Commit Name not found!")
		resText["err"] = "Error: nothing to commit"
		returnJs(w, resText)
	} else {
		tx.Commit()
		resText["err"] = ""
		returnJs(w, resText)
	}
}

func TxRollback(w http.ResponseWriter, r *http.Request) {
	rollbackName := r.URL.Query().Get("rbn")
	resText      := map[string]string{}
	
	fmt.Printf("[%-18s] Rollback name : %s\n", "TxRollback",rollbackName)
	fmt.Println("[* END Transaction *]");

	tx, ok := dbCommits[rollbackName]
	if !ok {
		fmt.Printf("[%-18s] Rollback name NOT FOUND: %s\n", "TxRollback",rollbackName)
		resText["err"] = "Error: nothing to rollback"
		returnJs(w, resText)
	} else {
		tx.Rollback()
		if rollbackName == "CompletePickList" {
			maxId := 0
			err := db.QueryRow("SELECT MAX(LID) from last_updated").Scan(&maxId)
			if err != nil {
				fmt.Println("ERR : ", err)
			}
			fmt.Printf("[%-18s] Last LID : %d in last_updated\n", "  Max LID", maxId)
			stmt := fmt.Sprintf("ALTER TABLE last_updated AUTO_INCREMENT = %d", maxId)
			res, err := db.Exec(stmt)
			if err != nil {
				fmt.Println("[ *--- AUTO_INCREMENT FAILED! ---*]\n", err)
			} else {
				afw, _ := res.RowsAffected()
				fmt.Println(" --- row affected:", afw)
			}
		}
		resText["err"] = ""
		returnJs(w, resText)
	}
}

func MessagePage(w http.ResponseWriter, r *http.Request) {
	tmplpath := "templates/tenda/messages.html"
	tmpl, err := template.ParseFiles(tmplpath)
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}