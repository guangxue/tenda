package tenda

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/guangxue/webapps/mysql"
)

func Stock(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	}

	if r.Method == "POST" {
	}

	if r.Method == "PUT" {
        SID := strings.TrimPrefix(r.URL.Path, "/tenda/api/stock/SID/")
        if strings.Contains(SID, "/") {
            SID = ""
        }
       err := r.ParseForm()
       if err != nil {
               fmt.Println("Form parse error:", err)
       }
        location := r.FormValue("location")
       model := r.FormValue("model")
       unit := r.FormValue("unit")
       cartons := r.FormValue("cartons")
       boxes := r.FormValue("boxes")
       total := r.FormValue("total")
       update_comments := r.FormValue("update_comments")

       if update_comments != "" {
           tx, ctx := mysql.Begin(db)
           updateInfo := map[string]interface{} {
                   "location":location,
                   "model":model,
                   "unit":unit,
                   "cartons":cartons,
                   "boxes":boxes,
                   "total":total,
                   "update_comments":update_comments,
           }
           fmt.Println("...UPDATE..ing stock_updated")
           mysql.
               Update(tbname["stock_updated"],false).
               Set(updateInfo).
               Where("SID", SID).
           With(tx, ctx)

           updatedStock := mysql.
               Select("SID", "location", "model", "unit", "cartons", "boxes", "total","update_comments").
               From("stock_updated").
               Where("SID", SID).
           With(tx, ctx)
           dbCommits["StockUpdate"] = tx
           returnJson(w, updatedStock)
       }
	}
}
