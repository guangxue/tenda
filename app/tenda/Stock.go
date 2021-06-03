package tenda

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/guangxue/webapps/mysql"
)

func Stock(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        allstock := mysql.
        Select("SID", "location", "model", "unit", "cartons", "boxes","total", "kind", "notes").
            From(tbname["stock_updated"]).
        Use(db)
        returnJson(w, allstock)
	}

	if r.Method == "POST" {
        err := r.ParseForm()
        if err != nil {
            fmt.Println("Form parse error:", err)
        }
        location := r.FormValue("location")
        model    := r.FormValue("model")
        unit     := r.FormValue("unit")
        cartons  := r.FormValue("cartons")
        boxes    := r.FormValue("boxes")
        total    := r.FormValue("total")
        kind     := r.FormValue("kind")

		insertStmt := map[string]interface{}{
			"location":location,
			"model":model,
			"unit":unit,
			"cartons":cartons,
			"boxes":boxes,
			"total":total,
            "kind":kind,
		}
        fmt.Println("InsertStatmentMap:", insertStmt)
        tx, ctx := mysql.Begin(db)
		insertFeedback := mysql.
            InsertInto(tbname["stock_updated"], insertStmt).
            With(tx, ctx)
        dbCommits["StockAdd"] = tx
        lastId, _ := insertFeedback[0]["lastId"]
        if lastId != "" {
            insertColumn := mysql.
                Select("SID", "location", "model", "unit", "cartons", "boxes", "total", "kind").
                From(tbname["stock_updated"]).
                Where("SID", lastId).
            With(tx,ctx)
            returnJson(w, insertColumn)
        }
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
        kind     := r.FormValue("kind")
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
                "kind": kind,
                "update_comments":update_comments,
            }
            fmt.Println("...UPDATE..ing stock_updated")
            mysql.
                Update(tbname["stock_updated"],false).
                Set(updateInfo).
                Where("SID", SID).
            With(tx, ctx)

            getUpdatedStock := mysql.
                Select("SID", "location", "model", "unit", "cartons", "boxes", "total","kind","update_comments").
                From(tbname["stock_updated"]).
                Where("SID", SID).
            With(tx, ctx)
            fmt.Println("Updatd stock in ctx:", getUpdatedStock)
            dbCommits["StockUpdate"] = tx
            returnJson(w, getUpdatedStock)
        }
    }
}
