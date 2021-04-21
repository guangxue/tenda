package tenda

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"html/template"
	"github.com/guangxue/webpages/apps/tenda/stock"
)

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
	querymodel := r.URL.Query().Get("model");
	queryall := r.URL.Query().Get("allmodels")
	querylocation := r.URL.Query().Get("location");
	fmt.Println("Request Path:", r.URL.Path)
	fmt.Println("querylocation:", querylocation)
	fmt.Println("querymodel:", querymodel)

    // get all models
	if len(queryall) > 0 {
        models := stock.GetAllModels();
	    ModelsJSON, err := json.Marshal(models)
	    if err != nil {
	    	fmt.Println("ModelsJson error: ", err)
	    }
		w.Write(ModelsJSON)
	}

    // get one model
	if len(querymodel) > 0 {
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

func ProcessForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("request Method:",r.Method)
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Form parse error:", err)
	}

	tableName := r.FormValue("tableName")
	if tableName == "picked" {
		PNO := r.FormValue("PNO")
		model := r.FormValue("model")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		now := r.FormValue("now")

		stock.InsertPicked(PNO, model, qty, customer, now)
	}
	
}

func TodaysPackages(w http.ResponseWriter, r *http.Request) {
	timenow := time.Now()
	currTime := timenow.Format("2006-01-02")
	fmt.Println("YYYY-MM-DD : ", currTime)

	allPicked := stock.GetTodayPackages(currTime)
	fmt.Printf("allPicked: %T",allPicked)
}