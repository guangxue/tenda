package tenda

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/guangxue/webapps/app/tenda/stock"
)

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
func QueryModels(w http.ResponseWriter, r *http.Request) {

    // Set Header for json HTTP response
	w.Header().Set("Content-Type", "application/json")

    // model name from URL querystring
    // -- api/models
    // -- api/models?model=AC18
    // -- api/models?location=0-G-1
	queryModel    := r.URL.Query().Get("model");
	queryLocation := r.URL.Query().Get("location");

	fmt.Println("[QueryModels] Request Path:", r.URL.Path)
	fmt.Println("[QueryModels] query Model:", queryModel, len(queryModel))
	fmt.Println("[QueryModels] query Location:", queryLocation, len(queryLocation))


    // get all models
	if len(queryModel) == 0 && len(queryLocation) == 0{
        models := stock.GetAllModels();
	    ModelsJSON, err := json.Marshal(models)
	    if err != nil {
	    	fmt.Println("ModelsJson error: ", err)
	    }
		w.Write(ModelsJSON)
	}

    // get one model
	if len(queryModel) > 0 {
        allModels := stock.GetModel(queryModel)
	    ModelsJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)

        w.Write(ModelsJSON)
	}

	// get location data
	if len(queryLocation) > 0 {
		allModels := stock.GetLocationModels(queryLocation);
		LocationJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
	    fmt.Println("LocationJSON", string(LocationJSON))
        w.Write(LocationJSON)
	}
}



func QueryLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryModel := r.URL.Query().Get("model");
	if len(queryModel) > 0 {
		allLocations := stock.GetModelLocations(queryModel);
		LocationJSON, err := json.Marshal(allLocations)
	    ErrorCheck(err)
	    fmt.Println("[QueryLocations] LocationJSON", string(LocationJSON))
        w.Write(LocationJSON)
	}
}


func PickedParcels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		date := r.URL.Query().Get("date");
		if date == "today" {
			timenow := time.Now()
			currTime := timenow.Format("2006-01-02")
			fmt.Println("YYYY-MM-DD : ", currTime)

			allPicked := stock.GetTodayPackages(currTime)
			PickedJSON, err := json.Marshal(allPicked)
		    if err != nil {
		    	fmt.Println("PickedJson error: ", err)
		    }
			w.Write(PickedJSON)
		}
		if date == "pending" {
			pendingParcels := stock.GetPendingParcels()
			ParcelJSON, err := json.Marshal(pendingParcels)
		    if err != nil {
		    	fmt.Println("ParcelJSON error: ", err)
		    }
			w.Write(ParcelJSON)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		PNO := r.FormValue("PNO")
		model := r.FormValue("model")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("location")
		now := r.FormValue("now")
		status := "Pending"
		lastId := stock.InsertPicked(PNO, model, qty, customer, location, status, now)
		insertResp := InsertResponse {lastId}
		resJSON, err := json.Marshal(insertResp)
	    if err != nil {
	    	fmt.Println("resJSON error: ", err)
	    }
		w.Write(resJSON)
	}
}

func QueryPickedWithPID (w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.PATH :", r.URL.Path)
}
