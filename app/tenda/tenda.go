package tenda

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	_ "strings"
	"encoding/json"
	_ "time"
	"strconv"
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

	queryModel    := r.URL.Query().Get("model");
	queryLocation := r.URL.Query().Get("location");

	fmt.Println("===> Request Path:", r.URL.Path)
	fmt.Println("=> query Model:", queryModel)
	fmt.Println("=> query Location:", queryLocation)


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
	    fmt.Printf("=> Locations\n")
	    for _, val := range allLocations {
	    	fmt.Printf("=> %s\n", val)
	    }
 	    w.Write(LocationJSON)
	}
}


func PickList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		date := r.URL.Query().Get("date");
		date = date + "%"
		status := r.URL.Query().Get("status");
		allPicked := mysql.Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").From("picklist").WhereLike("created_at",date).AndWhere("status", status).Use(db)
		PickedJSON, err := json.Marshal(allPicked)
	    if err != nil {
	    	fmt.Println("PickedJson error: ", err)
	    }
		w.Write(PickedJSON)
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
		status := "Pending"
		insertColumns := []string{"PNO","model","qty","customer","location","status"}
		insertValues  := []interface{}{PNO,model,qty,customer,location,status}
		insertFeedback := mysql.InsertInto("picklist", insertColumns, insertValues).Use(db)
		
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
	PID      int
	Model    string
	Qty      int
	Location string
}
type updateModel struct {
	PID      int
	Model    string
	Location string
	Cartons  int
	Boxes    int
	Total    int
}
func CompletePickList (w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("[CompletePicked] Form parse error:", err)
		}
		pickDate := r.FormValue("pickDate")
		pickStatus := r.FormValue("pickStatus")
		
		fmt.Println("==> [Complete Pick List]")
		fmt.Println(">> Pick data   :", pickDate)
		fmt.Println(">> pick status :", pickStatus)

		p := pickcolumns{}
		pcols := []pickcolumns{}
		if len(pickDate) == 0 {
			fmt.Println(">> [return] Empty pickDate NO ACTION made.")
			return
		}
		sqlstmt := "SELECT PID, model, qty, location FROM picklist WHERE created_at LIKE '"+pickDate+"%' AND status ='"+pickStatus+"'"
		fmt.Printf(" => :\n%s\n", sqlstmt)
		rows, err := db.Query(sqlstmt)
		if err != nil {
			fmt.Println("[CompletePicked] selection error:", err)
		}

		for rows.Next() {
			err := rows.Scan(&p.PID, &p.Model, &p.Qty, &p.Location)
			if err != nil {
				fmt.Println("[tenda.go CompletePicked]: dbColumns Scan error:207:", err)
			}
			pcols = append(pcols, p)
		}
		//-------------- Complete db.Query -------------/
		upModels := []updateModel{}
		for _, p := range pcols {
			fmt.Println("==> current pick(p) from picklist:")
			fmt.Println("{Model, Qty, Location}")
			fmt.Println(p)
			fmt.Println("==>")
			fmt.Println(">> Model:", p.Model)
			totals := 0
			unit := 0
			err := db.QueryRow("SELECT total,unit FROM stock_update WHERE model=? AND location=?", p.Model, p.Location).Scan(&totals, &unit)
			switch {
				case err == sql.ErrNoRows:
					fmt.Printf("[db *ERR*] no `totals` for model %s\n", p.Model)
				case err != nil:
					fmt.Printf("[db *ERR*] query error: %v\n", err)
				default:
					fmt.Printf(">> totals are %d\n", totals)
			}
			fmt.Println(">> pick qty:", p.Qty)
			newTotal := totals-p.Qty
			fmt.Println(">> *NEW Total:", newTotal)
			fmt.Printf(">> *unit are %d\n", unit)
			newCartons := newTotal/unit
			fmt.Println(">> *NEW Cartons:", newCartons)
			newBoxesFrac := float64(newTotal)/float64(unit) - float64(newCartons)
			fmt.Println(">> *NEW BoxeFrac:", newBoxesFrac)
			newBoxesFrac = newBoxesFrac * float64(unit)
			newBoxes := int(math.Round(newBoxesFrac))
			fmt.Println(">> *NEW Boxes:", newBoxes)
			upModel := updateModel{p.PID, p.Model, p.Location, newCartons, newBoxes, newTotal}
			upModels = append(upModels, upModel)
		}
		for _, val := range upModels {
			fmt.Printf("---\nUPDATE COLUMNS : [Model, Location, Cartons, Boxes, Total]\n")
			fmt.Printf("  SET  VALUES  : %v\n", val)
			fmt.Println("val:PID:", val.PID)
			fmt.Println("val:Model:", val.Model)
			fmt.Println("val:Location:", val.Location)
			fmt.Println("val:Cartons:", val.Cartons)
			fmt.Println("val:Boxes:", val.Boxes)
			fmt.Println("val:Total:", val.Total)

			updateSUInfo := map[string]interface{} {
				"cartons":val.Cartons,
				"boxes":val.Boxes,
				"total":val.Total,
			}
			updatePLInfo := map[string]interface{} {"status":"Complete"}
			mysql.Update("stock_update").Set(updateSUInfo).Where("model", val.Model).AndWhere("location", val.Location).Use(db)
			mysql.Update("picklist").Set(updatePLInfo).Where("PID", strconv.Itoa(val.PID)).Use(db)
		}
	}
}


func UpdatePickList(w http.ResponseWriter, r *http.Request) {

	dbPickedInfo := map[string]string{}

	if r.Method == "GET" {
		queryPID := r.URL.Query().Get("PID");
		fmt.Printf("==> Update Pick List (PID:%s)\n", queryPID)
		currentPID := mysql.Select("PNO", "model", "qty", "customer", "location", "status").From("picklist").Where("PID", queryPID).Use(db)
		dbPickedInfo = currentPID[0]
		dbPickedInfo["PID"] = queryPID
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
		fmt.Printf(">> PNO:%s\nmodel:%s\nqty:%s\ncustomer:%s\nlocation:%s\nstatus:%s",PNO, model, qty, customer,location,status)
		fmt.Println(">> dbPickedInfo::", dbPickedInfo)
		
		updateInfo := map[string]interface{} {
			"PNO":PNO,
			"model":model,
			"qty":qty,
			"customer":customer,
			"location":location,
			"status":status,
		}
		mysql.Update("picklist").Set(updateInfo).Where("PID",PID).Use(db)
	}
}