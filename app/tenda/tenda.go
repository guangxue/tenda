package tenda

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/guangxue/webapps/mysql"
)
var db = mysql.Connect("tenda");
type InsertResponse struct {
	LastId int64 `json:"lastId"`
}

func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tenda Index URL.Path->", r.URL.Path)
    fmt.Fprintf(w, "Tenda Pick and Pack System.")
}


func UpdatePickedPage(w http.ResponseWriter, r *http.Request) {
	queryPID := r.URL.Query().Get("PID");
	fmt.Println("[UpdatePickedPage] PID:", queryPID)
}
func Render(templateName string) http.HandlerFunc {
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

//----------------------------------------------------------*/
/*--------------------API-----------------------------------*/
func Models(w http.ResponseWriter, r *http.Request) {

    // Set Header for json HTTP response
	w.Header().Set("Content-Type", "application/json")

    // model name from URL querystring
    // -- api/models
    // -- api/models?model=AC18
    // -- api/models?location=0-G-1


    /********* TEST MySQL Statements **********
    // mysql.Select("model", "location").From("stock_locations").Where("model", "MW6-3PK").Use(db)
    // updateColumns := map[string]string{
    // 	"model": "ABC9-2PK",
    // 	"location":"0-0-0",
    // }
    // mysql.Update("stock_locations").Set(updateColumns).Where("SID", "12").Use(db)

    // mysql.Select("model","qty", "location").From("picked").WhereLike("last_updated", "2021-04-23%").Use(db)
    // insertColumns := []string{"PNO", "model", "qty", "customer", "location"}
    // insertValues  := []string{"PO20210412", "AC6", "eBay", "0-G-3"}
    // mysql.Insert("picked", insertColumns, insertValues)
    *******************************************/




	queryModel    := r.URL.Query().Get("model");
	queryLocation := r.URL.Query().Get("location");

	fmt.Println("[Models] Request Path:", r.URL.Path)
	fmt.Println("[Models] query Model:", queryModel)
	fmt.Println("[Models] query Location:", queryLocation)


    // get all models
	if len(queryModel) == 0 && len(queryLocation) == 0{
        modelNames := mysql.SelectDistinct("model").From("stock").Use(db);
	    ModelNamesJSON, err := json.Marshal(modelNames)
	    if err != nil {
	    	fmt.Println("ModelsJson error: ", err)
	    }
		w.Write(ModelNamesJSON)
	}

    // get one model
	if len(queryModel) > 0 {
        allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From("stock").Where("model", queryModel).Use(db)
	    ModelsJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
        w.Write(ModelsJSON)
	}

	// get location data
	if len(queryLocation) > 0 {
		allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From("stock").Where("location", queryLocation).Use(db)
		LocationJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
	    fmt.Println("LocationJSON", string(LocationJSON))
 	    w.Write(LocationJSON)
	}
}


func Locations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryModel := r.URL.Query().Get("model");
	
	if len(queryModel) > 0 {
		allLocations := mysql.Select("location").From("stock").Where("model", queryModel).Use(db)
		LocationJSON, err := json.Marshal(allLocations)
	    ErrorCheck(err)
	    fmt.Printf("[Locations]\nLocationJSON:%s\n", string(LocationJSON))
 	    w.Write(LocationJSON)
	}
}


func Picked(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		date := r.URL.Query().Get("date");
		if date == "today" {
			timenow := time.Now()
			timePattern := timenow.Format("2006-01-02")+"%"
			allPicked := mysql.Select("PID", "PNO", "model", "qty", "customer", "location", "status", "last_updated").From("picked").WhereLike("last_updated",timePattern).Use(db)
			PickedJSON, err := json.Marshal(allPicked)
		    if err != nil {
		    	fmt.Println("PickedJson error: ", err)
		    }
			w.Write(PickedJSON)
		}
		if date == "pending" {
			pendingParcels := mysql.Select("PID", "PNO", "model", "qty", "customer", "location", "status", "last_updated").From("picked").Where("status","Pending").Use(db)
			ParcelJSON, err := json.Marshal(pendingParcels)
		    if err != nil {
		    	fmt.Println("ParcelJSON error: ", err)
		    }
			w.Write(ParcelJSON)
		}
	}
	// Inserting picking informations

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		PNO := r.FormValue("PNO")
		model := r.FormValue("modelName")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("pickLocation")
		now := r.FormValue("now")
		status := "Pending"
		insertColumns := []string{"PNO","model","qty","customer","location","status", "last_updated"}
		insertValues  := []interface{}{PNO,model,qty,customer,location,status,now}
		mysql.InsertInto("picked", insertColumns, insertValues).Use(db)
		// insertResp := InsertResponse {lastId}
		// resJSON, err := json.Marshal(insertResp)
	 //    if err != nil {
	 //    	fmt.Println("resJSON error: ", err)
	 //    }
		// w.Write(resJSON)
	}
}

func QueryPickedWithPID (w http.ResponseWriter, r *http.Request) {
	// fmt.Println("r.PATH :", r.URL.Path)
}
