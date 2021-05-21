package tenda

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webapps/mysql"
)

func StockUpdate(w http.ResponseWriter, r *http.Request) {
	SID := r.URL.Query().Get("SID")
	tbName := r.URL.Query().Get("tbname")

	if SID != "" && r.Method == http.MethodGet {
		currentStockToUpdate := mysql.Select("SID", "location", "model", "unit", "cartons", "boxes","total", "update_comments").From(tbName).Where("SID", SID).Use(db);
		fmt.Println("currentStockToUpdate:", currentStockToUpdate)
		render(w, "stockupdate.html", currentStockToUpdate[0])
	}
}
