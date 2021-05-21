package tenda

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/guangxue/webapps/mysql"
)

func PickList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		date := r.URL.Query().Get("date");
		status := r.URL.Query().Get("status");

		model := strings.TrimPrefix(r.URL.Path, "/tenda/api/picklist/model/")
		if strings.Contains(model, "/") {
			model = ""
		}
		PID := strings.TrimPrefix(r.URL.Path, "/tenda/api/picklist/PID/")
		if strings.Contains(PID, "/") {
			PID = ""
		}
		PNO := strings.TrimPrefix(r.URL.Path, "/tenda/api/picklist/PNO/")
		if strings.Contains(PNO, "/") {
			PNO = ""
		}
		searchPNO := strings.TrimPrefix(r.URL.Path, "/tenda/api/picklist/search/PNO/")
		if strings.Contains(searchPNO, "/") {
			searchPNO = ""
		}

		fmt.Printf("[%-18s] path  :%s\n", "PickList", r.URL.Path);
		fmt.Printf("[%-18s] date  :%s\n", "PickList", date);
		fmt.Printf("[%-18s] status:%s\n", "PickList", status);
		fmt.Printf("[%-18s] PID   :%s\n", "PickList", PID);
		fmt.Printf("[%-18s] PNO   :%s\n", "PickList", PNO);
		fmt.Printf("[%-18s] searchPNO   :%s\n", "PickList", searchPNO);


		if status == "weeklycompleted" {
			// Weekly completed orders
			startDate := fmt.Sprintf("'%s'", date)
			endDate := fmt.Sprintf("date_add('%s', INTERVAL 7 DAY)", date)
			allPicked := mysql.
				Select("LID", "location", "model", "unit", "cartons", "boxes", "total", "completed_at").
				From(tbname["last_updated"]).
				WhereBetween("completed_at", startDate, endDate).
			Use(db)
			returnJson(w, allPicked)
		} else if PID != "" && status != "" {
			// Picked orders
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("PID",PID).
				AndWhere("status", "=",status).
			Use(db)
			returnJson(w, allPicked)
		} else if date != "" && status == "weeklypicked" { 
			// Weekly picked
			stmt := fmt.Sprintf("WITH weeklypicked AS (select model, qty, customer, location, status, created_at FROM picklist Where created_at BETWEEN '%s' AND date_add('%s', INTERVAL 7 DAY)) SELECT model, SUM(qty) as total from weeklypicked group by model", date, date)
			allPicked := mysql.
				SelectRaw(stmt, "model", "total").
			Use(db)
			fmt.Printf("[%-18s] PID   :%s %v\n", "weeklypicked:allpicked:", allPicked)
			returnJson(w, allPicked)
		} else if PNO != "" {
			// PNO
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("PNO",PNO).
			Use(db)
			returnJson(w, allPicked)
		} else if searchPNO != "" {
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				WhereLike("PNO","%"+searchPNO+"%").
			Use(db)
			returnJson(w, allPicked)
		} else if model != "" {
			if status == "from" && date != "" {
				allPicked := mysql.
					Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
					From(tbname["picklist"]).
					Where("model",model).
					AndWhere("created_at", ">", date).
				Use(db)
				returnJson(w, allPicked)
			} else {
				allPicked := mysql.
					Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
					From(tbname["picklist"]).
					Where("model",model).
				Use(db)
				returnJson(w, allPicked)
			}
		} else {
			odate := ""
			if date == "" {
				odate = "%"
			} else {
				odate = date + "%"
			}
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				WhereLike("created_at",odate).
				AndWhere("status", "=",status).
			Use(db)
			returnJson(w, allPicked)
		}
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		PNO := r.FormValue("PNO")
		model := r.FormValue("modelName")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("pickLocation")
		status := "Pending"
		insertQuery := map[string]interface{}{
			"PNO":PNO,
			"model":model,
			"qty":qty,
			"customer":customer,
			"location":location,
			"status":status,
		}
		insertFeedback := mysql.InsertInto(tbname["picklist"], insertQuery).Use(db)
		if err != nil {
			fmt.Println("Commit error:", err)
		}
		returnJson(w, insertFeedback)
	}

	if r.Method == "DELETE" {
		fmt.Println("......[] Recive delete method")
		PID := strings.TrimPrefix(r.URL.Path, "/tenda/api/picklist/")
		status := r.URL.Query().Get("status");
		fmt.Println("......[] DELETE PID :", PID)
		fmt.Println("......[] DELETE status :", status)
		
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

	if r.Method == "PUT" {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("Form parse error:", err)
		}
		PID := strings.TrimPrefix(r.URL.Path, "/tenda/api/picklist/PID/")
		if strings.Contains(PID, "/") {
			PID = ""
		}
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
