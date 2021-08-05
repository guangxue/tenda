package tenda

import (
	"fmt"
	"net/http"
    "strings"
	_ "github.com/guangxue/webapps/mysql"
)

func LastUpdated(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		LID := strings.TrimPrefix(r.URL.Path, "/tenda/api/lastupdated/LID/")
		if strings.Contains(LID, "/") {
			LID = ""
		}
		fmt.Println("[last updated API] LID:", LID)
		
		// location := r.FormValue("location")
		
	// 	tx, ctx := mysql.Begin(db)
	// 	updateInfo := map[string]interface{} {
	// 		"PNO":PNO,
	// 		"model":model,
	// 		"qty":qty,
	// 		"customer":customer,
	// 		"location":location,
	// 		"status":status,
	// 	}
	// 	mysql.Update(tbname["picklist"], false).Set(updateInfo).Where("PID",PID).With(tx, ctx)
	// 	fmt.Printf("[%-18s] Getting PID=%s from UPDATEd\n", "PickListUpdate", PID);
	// 	updatedPicked := mysql.
	// 		Select("PNO", "model", "qty", "customer", "location", "status").
	// 		From(tbname["picklist"]).
	// 		Where("PID", PID).
	// 		With(tx, ctx)
	// 	dbCommits["PickList"] = tx
	// 	fmt.Println("updatedPicked:", updatedPicked)
	// 	returnJson(w, updatedPicked)
	}
}
