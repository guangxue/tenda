package main

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webpages/apps/tenda"
	"github.com/guangxue/webpages/apps/blog"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    fmt.Fprintf(w, "gzhang webpages HomePage")
}

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/", HomePage)

	/**************************************************************/
	// Tenda pick and pack system
	mux.HandleFunc("/tenda", tenda.Index)
	mux.HandleFunc("/tenda/pick", tenda.PickPage)
	mux.HandleFunc("/tenda/query", tenda.QueryPage)
	mux.HandleFunc("/tenda/query/models", tenda.QueryModels)
	

	/**************************************************************/
	// Blog system
	mux.HandleFunc("/blog", blog.Admin)


	fmt.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Port listening error: ", err)
    }
}
