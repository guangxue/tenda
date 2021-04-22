package main

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webpages/apps/tenda"
	"github.com/guangxue/webpages/apps/blog"
	"github.com/guangxue/webpages/apps/home"
)

func main() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	/*------------------------------------------------------------*/
	// Main home page
	mux.HandleFunc("/", home.Index)


	/*------------------------------------------------------------*/
	// Tenda pick and pack system
	mux.HandleFunc("/tenda", tenda.Index)
	mux.HandleFunc("/tenda/pick", tenda.PickPage)
	mux.HandleFunc("/tenda/pack/",tenda.PackPage)
	mux.HandleFunc("/tenda/query", tenda.QueryPage)
	mux.HandleFunc("/tenda/update", tenda.UpdatePage)
	mux.HandleFunc("/tenda/api/models", tenda.QueryModels)
	mux.HandleFunc("/tenda/api/form-data",tenda.ProcessForm)
	mux.HandleFunc("/tenda/api/picked",tenda.TodaysPackages)
	

	/*------------------------------------------------------------*/
	// Blog system
	mux.HandleFunc("/blog", blog.Admin)


	/*------------------------------------------------------------*/
	fmt.Println("Listening on port :8080")
	err := http.ListenAndServe(":8080", mux)
    if err != nil {
        fmt.Println("Port listening error: ", err)
    }
}
