package tenda

import (
	"database/sql"
	"fmt"
	"time"
	"math"
	"net/http"
	_ "strings"
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

func render(w http.ResponseWriter, templateName string, data interface{}) {
	tmplpath := "templates/tenda/" + templateName
	tmpl, err := template.ParseFiles(tmplpath, "templates/tenda/base.html","templates/tenda/nav.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.Execute(w, data)
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

//----------------------------------------------------------*/
/*--------------------API-----------------------------------*/
func Models(w http.ResponseWriter, r *http.Request) {

    // Set Header for json HTTP response
	w.Header().Set("Content-Type", "application/json")

	/**********************************
    // model name from URL querystring
    // -- api/models
    // -- api/models?model=AC18
    // -- api/models?location=0-G-1
    ***********************************/


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
        modelNames := mysql.SelectDistinct("model").From("stock_update").Use(db);
	    ModelNamesJSON, err := json.Marshal(modelNames)
	    if err != nil {
	    	fmt.Println("ModelsJson error: ", err)
	    }
		w.Write(ModelNamesJSON)
	}

    // get one model
	if len(queryModel) > 0 {
        allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From("stock_update").Where("model", queryModel).Use(db)
	    ModelsJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
        w.Write(ModelsJSON)
	}

	// get location data
	if len(queryLocation) > 0 {
		allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From("stock_update").Where("location", queryLocation).Use(db)
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
		allLocations := mysql.Select("location").From("stock_update").Where("model", queryModel).Use(db)
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
		fmt.Println("now", now)
		status := "Pending"
		insertColumns := []string{"PNO","model","qty","customer","location","status", "last_updated"}
		insertValues  := []interface{}{PNO,model,qty,customer,location,status,now}
		insertFeedback := mysql.InsertInto("picked", insertColumns, insertValues).Use(db)
		
		respJSON, err := json.Marshal(insertFeedback)
	    if err != nil {
	    	fmt.Println("respJSON error: ", err)
	    }
		w.Write(respJSON)
	}
}

func QueryPickedWithPID (w http.ResponseWriter, r *http.Request) {
	// fmt.Println("r.PATH :", r.URL.Path)
}


type pickcolumns struct {
	Model    string
	Qty      int
	Location string
}
type updateModel struct {
	Model    string
	Location string
	Cartons  int
	Boxes    int
	Total    int
}
func CompletePicked (w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("[CompletePicked] Form parse error:", err)
		}
		completeType := r.FormValue("completeType")
		fmt.Println("completeType:", completeType)

		p := pickcolumns{}
		pcols := []pickcolumns{}
		sqlstmt := "SELECT model, qty, location FROM picked "
		if completeType == "Today" {
			timenow := time.Now()
			timePattern := timenow.Format("2006-01-02")+"%"
			sqlstmt += "WHERE status='Pending' AND last_updated LIKE '"+timePattern+"'"
			fmt.Printf("[db] SELECTing ... picked\n%s\n",sqlstmt)
		}
		if completeType == "Pending" {
			sqlstmt += "WHERE status='Pending'"
			fmt.Printf("[db] SELECTing ... picked\n%s\n",sqlstmt)
		}
		if completeType == "Complete" {
			sqlstmt += "WHERE status='Complete'"
			fmt.Printf("[db] SELECTing ... picked\n%s\n",sqlstmt)
		}
		rows, err := db.Query(sqlstmt)
		if err != nil {
			fmt.Println("[CompletePicked] selection error:", err)
		}

		for rows.Next() {
			err := rows.Scan(&p.Model, &p.Qty, &p.Location)
			if err != nil {
				fmt.Println("[tenda.go CompletePicked]: dbColumns Scan error:207:", err)
			}
			pcols = append(pcols, p)
		}
		//-------------- Complete db.Query -------------/
		upModels := []updateModel{}
		for _, p := range pcols {
			fmt.Println("====> pcols:", p)
			fmt.Println("=> p.Model:", p.Model)
			totals := 0
			unit := 0
			err := db.QueryRow("SELECT total,unit FROM stock_update WHERE model=? AND location=?", p.Model, p.Location).Scan(&totals, &unit)
			switch {
				case err == sql.ErrNoRows:
					fmt.Printf("[db *ERR*] no `totals` for model %s\n", p.Model)
				case err != nil:
					fmt.Printf("[db *ERR*] query error: %v\n", err)
				default:
					fmt.Printf("=> totals are %d\n", totals)
			}
			fmt.Println("=> current pick:", p.Qty)
			newTotal := totals-p.Qty
			fmt.Println("=> *NEW Total:", newTotal)
			fmt.Printf("=> [unit are %d]\n", unit)
			newCartons := newTotal/unit
			fmt.Println("=> *NeW Cartons:", newCartons)
			newBoxesFrac := float64(newTotal)/float64(unit) - float64(newCartons)
			fmt.Println("=> *NEW BoxeFrac:", newBoxesFrac)
			newBoxesFrac = newBoxesFrac * float64(unit)
			newBoxes := int(math.Round(newBoxesFrac))
			fmt.Println("=> *NEW Boxes:", newBoxes)
			upModel := updateModel{p.Model, p.Location, newCartons, newBoxes, newTotal}
			upModels = append(upModels, upModel)
		}
		for _, val := range upModels {
			fmt.Println("Update Models :", val)
		}
	}
}


func UpdatePickedPage(w http.ResponseWriter, r *http.Request) {

	dbPickedInfo := map[string]string{}

	if r.Method == "GET" {
		queryPID := r.URL.Query().Get("PID");
		fmt.Println("[UpdatePickedPage] PID:", queryPID)
		currentPID := mysql.Select("PNO", "model", "qty", "customer", "location", "status").From("picked").Where("PID", queryPID).Use(db)
		dbPickedInfo = currentPID[0]
		dbPickedInfo["PID"] = queryPID
		fmt.Println("CurrentPID: info:", dbPickedInfo)
		render(w, "update_picked.html", dbPickedInfo)
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		PID := r.FormValue("PID")
		PNO := r.FormValue("PNO")
		model := r.FormValue("model")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("location")
		status := r.FormValue("statusoption")
		timenow := r.FormValue("timenow")
		fmt.Printf("PNO:%s\nmodel:%s\nqty:%s\ncustomer:%s\nlocation:%s\nstatus:%s",PNO, model, qty, customer,location,status)
		fmt.Println("dbPickedInfo::", dbPickedInfo)
		
		updateInfo := map[string]string {
			"PNO":PNO,
			"model":model,
			"qty":qty,
			"customer":customer,
			"location":location,
			"status":status,
			"last_updated":timenow,
		}
		mysql.Update("picked").Set(updateInfo).Where("PID",PID).Use(db)
	}
}