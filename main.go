package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	tenda "tenda/utils"
)

var mux = http.NewServeMux()

func main() {
	fs := http.FileServer(http.Dir("./static/tenda/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	/*------------------------------------------------------------*/
	// Tenda Path Routing
	mux.HandleFunc("/{$}", tenda.RenderHandler("search.html"))
	mux.HandleFunc("/packingslip", tenda.RenderHandler("packingslip.html"))
	mux.HandleFunc("/picklist", tenda.RenderHandler("picklist.html"))
	mux.HandleFunc("/update", tenda.RenderHandler("update.html"))
	mux.HandleFunc("/stock/add", tenda.RenderHandler("stockadd.html"))
	mux.HandleFunc("/stock/update", tenda.StockUpdatePage)
	mux.HandleFunc("/stock", tenda.RenderHandler("stock.html"))
	mux.HandleFunc("/soh", tenda.RenderHandler("soh.html"))
	mux.HandleFunc("/picklist/update", tenda.PickListUpdatePage)
	mux.HandleFunc("/picklist/inspect", tenda.PickListInspectPage)
	mux.HandleFunc("/yam", tenda.MessagePage)
	mux.HandleFunc("/lastupdated", tenda.LastUpdatedPage)

	// Tenda API system
	mux.HandleFunc("/api/model/", tenda.Model)
	mux.HandleFunc("/api/locations", tenda.Locations)
	mux.HandleFunc("/api/picklist/complete", tenda.PickListComplete)
	mux.HandleFunc("/api/picklist/", tenda.PickList)
	mux.HandleFunc("/api/lastupdated", tenda.LastUpdated)
	mux.HandleFunc("PUT /api/stock/SID/{sid}", tenda.PutStock)
	mux.HandleFunc("/api/soh", tenda.SOH)
	mux.HandleFunc("/api/txcm", tenda.TxCommit)
	mux.HandleFunc("/api/txrb", tenda.TxRollback)

	fmt.Printf("[%-18s] Listening HTTPS on port :8080\n", " -- main.go")
	TLSerr := http.ListenAndServeTLS(":8080", "/home/guangxue/ssl/guangxuezhang_com_chain.crt", "/home/guangxue/ssl/ecc.key", mux)
	if TLSerr != nil {
		fmt.Println("ListenAndServe: ", TLSerr)
	}
}

func rootRouting(w http.ResponseWriter, r *http.Request) {
	rPath := r.URL.Path
	fmt.Printf("[ **%-12s**:] %s\n", " Fetch URL<routing>", rPath)

	if strings.Compare(rPath, "/") == 0 {
		mux.HandleFunc("/", tenda.RenderHandler("search.html"))
	}

	if rPath != "/" {
		// render 404 page
		tmpl, err := template.ParseFiles("templates/404.html")
		templVar := map[string]interface{}{
			"path": r.URL.Path,
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

	// render HOME: search
	tmpl, err := template.ParseFiles("templates/search.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}
