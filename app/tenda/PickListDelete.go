package tenda

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webapps/mysql"
)

func PickListDelete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Form parse error:", err)
	}
	PID := r.FormValue("PID");
	status := r.FormValue("status");
	fmt.Println("PID:", PID)
	fmt.Println("status:",status)
	
	if status == "Pending" {
		// delete
		rowsAffected := mysql.DeleteFrom(tbname["picklist"], false).Where("PID", PID).Use(db)
		writeJson := map[string]string{
			"rowsAffected": rowsAffected[0]["rowsAffected"],
		}
		returnJs(w, writeJson)
	}

	if status == "Complete" {
		// rollback
	}
}