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
	mux.HandleFunc("/tenda", tenda.Login)
	mux.HandleFunc("/tenda/pick", tenda.RenderHandler("packingslip.html"))
	mux.HandleFunc("/tenda/picklist",tenda.RenderHandler("picklist.html"))
	mux.HandleFunc("/tenda/find", tenda.RenderHandler("find.html"))
	mux.HandleFunc("/tenda/stocktakes", tenda.Stocktakes)
	mux.HandleFunc("/tenda/update", tenda.RenderHandler("update.html"))
	mux.HandleFunc("/tenda/update/stock", tenda.UpdateStock)
	mux.HandleFunc("/tenda/update/picklist", tenda.UpdatePickList)
	mux.HandleFunc("/tenda/stock/add", tenda.AddStock)

	// Tenda API system
	mux.HandleFunc("/tenda/api/models", tenda.Models)
	mux.HandleFunc("/tenda/api/locations", tenda.Locations)
	mux.HandleFunc("/tenda/api/picklist",tenda.PickList)
	mux.HandleFunc("/tenda/api/picked/delete",tenda.PickedDelete)
	mux.HandleFunc("/tenda/api/complete/picklist",tenda.CompletePickList)

	/*------------------------------------------------------------*/
	// Blog system
	mux.HandleFunc("/blog", blog.Admin)


	/*------------------------------------------------------------*/
	fmt.Printf("[%-18s] Listening on port :8080\n", "Main")
	err := http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Port listening error: ", err)
    }
}


func routing(w http.ResponseWriter, r *http.Request) {
	rPath := r.URL.Path
	fmt.Printf("[%-18s] Request path: %s\n", "Main", rPath)

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
