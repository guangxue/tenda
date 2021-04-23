package blog

import (
    "net/http"
    "html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/blog" {
        http.NotFound(w, r)
        return
    }

    template_files := []string {
        "./templates/blog/home.html",
        "./templates/blog/base.html",
        "./templates/blog/footer.html",
    }
    ts, err := template.ParseFiles(template_files...)
    if err != nil {
        http.Error(w, err.Error(), 500)
    }
    err = ts.Execute(w,nil)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}


func Admin(w http.ResponseWriter, r *http.Request) {

    template_files := []string {
        "./templates/blog/admin/editor.html",
    }
    ts, err := template.ParseFiles(template_files...)
    if err != nil {
        http.Error(w, err.Error(), 500)
    }
    err = ts.Execute(w,nil)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
}