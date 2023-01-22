package main

import (
	"fmt"
	"net/http"
	"html/template"
	"strings"
	"github.com/guangxue/webapps/app/tenda"
)

var mux = http.NewServeMux()

func main() {
	fs := http.FileServer(http.Dir("./static/tenda/"))
	mux.Handle("/tenda/static/", http.StripPrefix("/tenda/static/", fs))

	mux.HandleFunc("/", routing)

	/*------------------------------------------------------------*/
	// Tenda pick and pack system
	mux.HandleFunc("/tenda/", tenda.Index)
	mux.HandleFunc("/tenda/packingslip", tenda.RenderHandler("packingslip.html"))
	mux.HandleFunc("/tenda/picklist",tenda.RenderHandler("picklist.html"))
	mux.HandleFunc("/tenda/search", tenda.RenderHandler("search.html"))
	mux.HandleFunc("/tenda/update", tenda.RenderHandler("update.html"))
	mux.HandleFunc("/tenda/stock/add", tenda.RenderHandler("stockadd.html"))
	mux.HandleFunc("/tenda/stock/update", tenda.StockUpdatePage)
	mux.HandleFunc("/tenda/stock", tenda.RenderHandler("stock.html"))
	mux.HandleFunc("/tenda/soh", tenda.RenderHandler("soh.html"))
	mux.HandleFunc("/tenda/picklist/update", tenda.PickListUpdatePage)
	mux.HandleFunc("/tenda/picklist/inspect", tenda.PickListInspectPage)
	mux.HandleFunc("/tenda/yam", tenda.MessagePage)
	mux.HandleFunc("/tenda/lastupdated", tenda.LastUpdatedPage)
	
	// Tenda API system
	mux.HandleFunc("/tenda/api/locations", tenda.Locations)
	mux.HandleFunc("/tenda/api/picklist/complete",tenda.PickListComplete)
	mux.HandleFunc("/tenda/api/lastupdated", tenda.LastUpdated)
	mux.HandleFunc("/tenda/api/txcm",tenda.TxCommit)
	mux.HandleFunc("/tenda/api/txrb",tenda.TxRollback)


	fmt.Printf("[%-18s] Listening HTTPS on port :8080\n", " -- main.go")
	// err := http.ListenAndServe(":8080", mux)
	// if err != nil {
	// 	fmt.Println("Port listening error: ", err)
	// }
 	TLSerr := http.ListenAndServeTLS(":8080", "/home/guangxue/ssl/guangxuezhang_com_chain.crt", "/home/guangxue/ssl/ecc.key", mux)
 	if TLSerr != nil {
 		fmt.Println("ListenAndServe: ", TLSerr)
 	}
}


func routing(w http.ResponseWriter, r *http.Request) {
	rPath := r.URL.Path
	fmt.Printf("[ **%-12s**:] %s\n", " Fetch URL<routing>", rPath)

	if strings.HasPrefix(rPath, "/tenda/api/model") {
		tenda.Model(w, r)
		return
	}

	if strings.HasPrefix(rPath, "/tenda/api/picklist") {
		tenda.PickList(w, r)
		return
	}

    if strings.HasPrefix(rPath, "/tenda/api/stock") {
        tenda.Stock(w, r)
        return
    }

    if strings.HasPrefix(rPath, "/tenda/api/soh") {
        tenda.SOH(w, r)
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
