package home

import (
	"fmt"
	"net/http"
	"html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
 //        http.NotFound(w, r)
 //        return
 //    }
	
	fmt.Println("Request path: ", r.URL.Path)
	if r.URL.Path != "/" {
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
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}