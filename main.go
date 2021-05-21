package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strings"
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
	mux.HandleFunc("/tenda/packingslip", tenda.RenderHandler("packingslip.html"))
	mux.HandleFunc("/tenda/picklist",tenda.RenderHandler("picklist.html"))
	mux.HandleFunc("/tenda/search", tenda.RenderHandler("search.html"))
	mux.HandleFunc("/tenda/update", tenda.RenderHandler("update.html"))
	mux.HandleFunc("/tenda/stock/add", tenda.RenderHandler("stockadd.html"))
	mux.HandleFunc("/tenda/stock/update", tenda.StockUpdate)
	mux.HandleFunc("/tenda/stocktakes", tenda.Stocktakes)
	mux.HandleFunc("/tenda/picklist/update", tenda.PickListUpdate)

	// Tenda API system
	mux.HandleFunc("/tenda/api/models", tenda.Models)
	mux.HandleFunc("/tenda/api/locations", tenda.Locations)
	mux.HandleFunc("/tenda/api/picklist/complete",tenda.PickListComplete)
	mux.HandleFunc("/tenda/api/txcm",tenda.TxCommit)
	mux.HandleFunc("/tenda/api/txrb",tenda.TxRollback)

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

	if strings.HasPrefix(rPath, "/tenda/api/picklist") {
		fmt.Println("start withs ::/tenda/api/picklist")
		tenda.PickList(w, r)
		return
	}

    if strings.HasPrefix(rPath, "/tenda/api/stock") {
        tenda.Stock(w, r)
        return
    }

	if rPath != "/" {
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
