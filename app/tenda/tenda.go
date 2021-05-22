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
	// tbname["stock_updated"] = "stock_updated"
	// tbname["last_updated"] = "last_updated"
	// tbname["picklist"] = "picklist"
	// tbname["stocktakes"] = "stocktakes"

	tbname["stock_updated"] = "stock_updated_test"
	tbname["last_updated"] = "last_updated_test"
	tbname["picklist"] = "picklist_test"
	tbname["stocktakes"] = "stocktakes"
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

func Models(w http.ResponseWriter, r *http.Request) {

	queryModel    := r.URL.Query().Get("model");
	queryLocation := r.URL.Query().Get("location");

	fmt.Printf("[%-18s] Request Path:%s\n", "Models", r.URL.Path)
	fmt.Printf("[%-18s] query Model:%s\n", "Models", queryModel)
	fmt.Printf("[%-18s] query Location:%s\n", "Models", queryLocation)


    // get all models
	if len(queryModel) == 0 && len(queryLocation) == 0{
        modelNames := mysql.SelectDistinct("model").From(tbname["stock_updated"]).Use(db);
        returnJson(w, modelNames)
	}

    // get one model
	if len(queryModel) > 0 {
        allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From(tbname["stock_updated"]).Where("model", queryModel).Use(db)
	    returnJson(w, allModels)
	}

	// get location data
	if len(queryLocation) > 0 {
		allModels := mysql.
			Select("model", "location", "unit", "cartons", "boxes", "total").
			From(tbname["stock_updated"]).
			Where("location", queryLocation).
		Use(db)
		returnJson(w, allModels)
	}
}

func Locations(w http.ResponseWriter, r *http.Request) {
	searchModel := r.URL.Query().Get("model");
	if len(searchModel) > 0 {
		allLocations := mysql.
			Select("location").
			From(tbname["stock_updated"]).
			Where("model", searchModel).
		Use(db)
		returnJson(w, allLocations)
	}
}

func PickListUpdatePage(w http.ResponseWriter, r *http.Request) {

	dbPickedInfo := map[string]string{}
	
	if r.Method == "GET" {
		queryPID := r.URL.Query().Get("PID")
		LID := r.URL.Query().Get("LID")

		if queryPID != "" {
			fmt.Printf("[%-18s] PID:%s\n", "PickListUpdate",queryPID)
			currentPID := mysql.Select("PNO", "model", "qty", "customer", "location", "status").From(tbname["picklist"]).Where("PID", queryPID).Use(db)
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

func StockUpdatePage(w http.ResponseWriter, r *http.Request) {
	SID := r.URL.Query().Get("SID")

	if SID != "" && r.Method == http.MethodGet {
		currentStockToUpdate := mysql.
			Select("SID", "location", "model", "unit", "cartons", "boxes","total", "update_comments").
			From(tbname["stock_updated"]).
			Where("SID", SID).
			Use(db)
		fmt.Println("currentStockToUpdate:", currentStockToUpdate)
		render(w, "stockupdate.html", currentStockToUpdate[0])
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
		resText["err"] = ""
		returnJs(w, resText)
	}
}
