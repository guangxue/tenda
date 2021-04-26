package main

import (
	"fmt"
	"net/http"
	"html/template"
	"github.com/guangxue/webapps/app/tenda"
	"github.com/guangxue/webapps/app/blog"
)

var mux = http.NewServeMux()

func main() {
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", routing)


	/*------------------------------------------------------------*/
	// Tenda pick and pack system
	mux.HandleFunc("/tenda", tenda.Index)
	mux.HandleFunc("/tenda/pick", tenda.Render("pick.html"))
	mux.HandleFunc("/tenda/picklist",tenda.Render("picklist.html"))
	mux.HandleFunc("/tenda/query", tenda.Render("query.html"))
	mux.HandleFunc("/tenda/update", tenda.Render("update.html"))
	mux.HandleFunc("/tenda/update/picked", tenda.UpdatePickedPage)

	// Tenda API system
	mux.HandleFunc("/tenda/api/models", tenda.Models)
	mux.HandleFunc("/tenda/api/locations", tenda.Locations)
	mux.HandleFunc("/tenda/api/picked",tenda.Picked)
	mux.HandleFunc("/tenda/api/picked/pid",tenda.QueryPickedWithPID)

	/*------------------------------------------------------------*/
	// Blog system
	mux.HandleFunc("/blog", blog.Admin)


	/*------------------------------------------------------------*/
	fmt.Println("[Main] Listening on port :8080")
	err := http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Port listening error: ", err)
    }
}


func routing(w http.ResponseWriter, r *http.Request) {
	rPath := r.URL.Path
	fmt.Println("[Main] Request path: ", rPath)

	if r.URL.Path != "/" {
		// render 404 page
		tmpl, err := template.ParseFiles("templates/404.html")
		templVar := map[string]interface{}{
			"path":r.URL.Path,
		}
		if err != nil {
			fmt.Println("template parsing errors: ", err)
		}
		err = tmpl.Execute(w, templVar)
		if err != nil {
			fmt.Println("template executing errors: ", err)
		}
		return
	}

	// render HOME 
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}
