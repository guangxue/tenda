package tenda

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webapps/mysql"
)

func PickListUpdate(w http.ResponseWriter, r *http.Request) {

	dbPickedInfo := map[string]string{}
	
	if r.Method == "GET" {
		queryPID := r.URL.Query().Get("PID")
		LID := r.URL.Query().Get("LID")

		if queryPID != "" {
			fmt.Printf("[%-18s] PID:%s\n", "PickListUpdate",queryPID)
			currentPID := mysql.Select("PNO", "model", "qty", "customer", "location", "status").From(tbname["picklist"]).Where("PID", queryPID).Use(db)
			dbPickedInfo = currentPID[0]
			dbPickedInfo["PID"] = queryPID
			dbPickedInfo["status2"] = ""
			if dbPickedInfo["status"] == "Pending" {
				dbPickedInfo["status2"] = "Complete"
			} else {
				dbPickedInfo["status2"] = "Pending"
			}
			fmt.Printf("[%-18s] /picklist GET PID:%s\n", "UpdatePickList", queryPID)
			fmt.Printf("[%-18s] return:%s\n", "UpdatePickList", dbPickedInfo)
			render(w, "picklistupdate.html", dbPickedInfo)
		}
		if LID != "" {
			currentLID := mysql.
				Select("LID", "location", "model", "unit", "cartons", "boxes", "total", "completed_at").
				From(tbname["last_updated"]).
				Where("LID", LID).
			Use(db)
			returnJson(w, currentLID)
		}
	}

	if r.Method == "PUT" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		PID := r.FormValue("PID")
		PNO := r.FormValue("PNO")
		model := r.FormValue("model")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("location")
		status := r.FormValue("status")
		
		tx, ctx := mysql.Begin(db)
		updateInfo := map[string]interface{} {
			"PNO":PNO,
			"model":model,
			"qty":qty,
			"customer":customer,
			"location":location,
			"status":status,
		}
		mysql.Update(tbname["picklist"], false).Set(updateInfo).Where("PID",PID).With(tx, ctx)
		fmt.Printf("[%-18s] Getting PID=%s from UPDATEd\n", "PickListUpdate", PID);
		updatedPicked := mysql.
			Select("PNO", "model", "qty", "customer", "location", "status").
			From(tbname["picklist"]).
			Where("PID", PID).
			With(tx, ctx)
		dbCommits["PickList"] = tx
		fmt.Println("updatedPicked:", updatedPicked)
		returnJson(w, updatedPicked)
	}
}