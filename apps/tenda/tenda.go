package tenda

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/guangxue/webpages/apps/tenda/stock"
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

func PickPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request path ->", r.URL.Path);
	tmpl, err := template.ParseFiles("templates/tenda/base.html", "templates/tenda/nav.html", "templates/tenda/pick.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}

func QueryPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/tenda/base.html", "templates/tenda/nav.html", "templates/tenda/query.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}
func UpdatePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/tenda/base.html", "templates/tenda/nav.html", "templates/tenda/update.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}

func QueryModels(w http.ResponseWriter, r *http.Request) {

    // Set Header for json HTTP response
	w.Header().Set("Content-Type", "application/json")

    // model name from URL querystring
	querymodel := r.URL.Query().Get("model"); // model=all OR model=ACU10
	querylocation := r.URL.Query().Get("location"); // location=0-G-1


	fmt.Println("[QueryModels] Request Path:", r.URL.Path)
	fmt.Println("[QueryModels] querylocation:", querylocation)
	fmt.Println("[QueryModels] querymodel:", querymodel)

    // get all models
	if querymodel == "all" {
        models := stock.GetAllModels();
	    ModelsJSON, err := json.Marshal(models)
	    if err != nil {
	    	fmt.Println("ModelsJson error: ", err)
	    }
		w.Write(ModelsJSON)
	}

    // get one model
	if len(querymodel) > 0 && querymodel != "all" {
        allModels := stock.GetModel(querymodel)
	    ModelsJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)

        w.Write(ModelsJSON)
	}

	// get location data
	if len(querylocation) > 0 {
		allModels := stock.GetLocationModels(querylocation);
		ModelsJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
	    fmt.Println("ModelsJSON %v\n", string(ModelsJSON))
        w.Write(ModelsJSON)
	}
}


func PickListPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/tenda/base.html", "templates/tenda/nav.html", "templates/tenda/picklist.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
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
		now := r.FormValue("now")

		lastId := stock.InsertPicked(PNO, model, qty, customer, now)
		insertResp := InsertResponse {lastId}
		resJSON, err := json.Marshal(insertResp)
	    if err != nil {
	    	fmt.Println("resJSON error: ", err)
	    }
		w.Write(resJSON)
	}
}