package tenda

import (
	"fmt"
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