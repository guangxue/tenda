package tenda

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webapps/mysql"
)

func PickList(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		date := r.URL.Query().Get("date");
		status := r.URL.Query().Get("status");
		PID := r.URL.Query().Get("PID");
		PNO := r.URL.Query().Get("PNO");
		model := r.URL.Query().Get("model");

		fmt.Printf("[%-18s] date  :%s\n", "PickList", date);
		fmt.Printf("[%-18s] status:%s\n", "PickList", status);
		fmt.Printf("[%-18s] PID   :%s\n", "PickList", PID);
		fmt.Printf("[%-18s] PNO   :%s\n", "PickList", PNO);

		if status == "completed_at" {
			// Weekly completed orders
			startDate := fmt.Sprintf("'%s'", date)
			endDate := fmt.Sprintf("date_add('%s', INTERVAL 7 DAY)", date)
			allPicked := mysql.
				Select("LID", "location", "model", "unit", "cartons", "boxes", "total", "completed_at").
				From(tbname["last_updated"]).
				WhereBetween("completed_at", startDate, endDate).
			Use(db)
			returnJson(w, allPicked)
		} else if len(PID) > 0 && len(status) > 0 {
			// Picked orders
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("PID",PID).
				AndWhere("status", "=",status).
			Use(db)
			returnJson(w, allPicked)
		} else if len(date) > 0 && status == "created_at" { 
			// Weekly picked
			stmt := fmt.Sprintf("WITH weeklypicked AS (select model, qty, customer, location, status, created_at FROM picklist Where created_at BETWEEN '%s' AND date_add('%s', INTERVAL 7 DAY)) SELECT model, SUM(qty) as total from weeklypicked group by model", date, date)
			allPicked := mysql.
				SelectRaw(stmt, "model", "total").
			Use(db)
			fmt.Printf("[%-18s] PID   :%s %v\n", "weeklypicked:allpicked:", allPicked)
			returnJson(w, allPicked)
		} else if len(PNO) > 0 {
			// PNO
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("PNO",PNO).
			Use(db)
			returnJson(w, allPicked)
		} else if model != "" {
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("model",model).
			Use(db)
			returnJson(w, allPicked)
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

	// Inserting picking informations
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
}

func PickedDelete(w http.ResponseWriter, r *http.Request) {
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
		tx,ctx := mysql.Begin(db)
		mysql.DeleteFrom(tbname["picklist"], false).Where("PID", PID).With(tx,ctx)
	}

	if status == "Complete" {
		// rollback
	}
}