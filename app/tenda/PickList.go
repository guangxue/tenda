package tenda

import (
	"fmt"
	"net/http"
	"strings"
	"time"
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
				Select("LID", "location", "model", "unit","old_total","total_picks","cartons", "boxes", "total", "completed_at").
				From(tbname["last_updated"]).
				WhereBetween("completed_at", startDate, endDate).
			Use(db)
			returnJson(w, allPicked)
		} else if PID != "" && status != "" {
			// Picked orders
			allPicked := mysql.
				Select("PID", "PNO","sales_mgr", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("PID",PID).
				AndWhere("status", "=",status).
			Use(db)
			returnJson(w, allPicked)
		} else if date != "" && status == "weeklypicked" {
			// Weekly picked
			stmt := fmt.Sprintf("select pno, customer, sales_mgr,model, qty, created_at FROM picklist Where created_at BETWEEN '%s' AND date_add('%s', INTERVAL 7 DAY)", date, date)
			allPicked := mysql.
				SelectRaw(stmt, "pno", "customer", "sales_mgr","model", "qty", "created_at").
			Use(db)
			if len(allPicked) > 0 {
				allPicked[0]["weeklypicked"] = "1"
			}
			fmt.Printf("[%-18s] PID   :%s %v\n", "weeklypicked:allpicked:", allPicked)
			returnJson(w, allPicked)
		} else if PNO != "" {
			// PNO
			allPicked := mysql.
				Select("PID", "PNO","sales_mgr", "model", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				Where("PNO",PNO).
			Use(db)
			returnJson(w, allPicked)
		} else if searchPNO != "" {
			allPicked := mysql.
				Select("PID", "PNO", "model","sales_mgr", "qty", "customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
				WhereLike("PNO","%"+searchPNO+"%").
			Use(db)
			returnJson(w, allPicked)
		} else if model != "" {
			if status == "from" && date != "" {
				allPicked := mysql.
					Select("PID", "PNO", "model", "qty","sales_mgr", "customer", "location", "status", "created_at", "updated_at").
					From(tbname["picklist"]).
					Where("model",model).
					AndWhere("created_at", ">", date).
				Use(db)
				returnJson(w, allPicked)
			} else {
				allPicked := mysql.
					Select("PID", "PNO", "model", "qty","sales_mgr","customer", "location", "status", "created_at", "updated_at").
					From(tbname["picklist"]).
					Where("model",model).
				Use(db)
				returnJson(w, allPicked)
			}
		} else if status == "all" {
			allPicked := mysql.
				Select("PID", "PNO", "model", "qty","sales_mgr","customer", "location", "status", "created_at", "updated_at").
				From(tbname["picklist"]).
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
				Select("PID", "PNO", "model","sales_mgr", "qty", "customer", "location", "status", "created_at", "updated_at").
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
		model := r.FormValue("model")
		qty := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("pickLocation")
		salesMgr := r.FormValue("sales_mgr")
		status := "Pending"
		td := time.Now()
		todayDate  := td.Format("20060102 15:04:05")
		today,tdTime  := strings.Split(todayDate, " ")[0], strings.Split(todayDate, " ")[1]
		fmt.Println("todayDate:", today)
		fmt.Println("tdTime:", tdTime)
		// AE20210816-1
		PNOdt := strings.Split(PNO, "-")[0]
		fmt.Println("PNOdt:", PNOdt)
		PNOdate := PNOdt[2:]
		fmt.Println("PNODate:", PNOdate)
		fmt.Println("PNOdate len:", len(PNOdate))

		if len(PNOdate) != 8 {
			resText := map[string]string {
				"err": "Wrong Date Format",
			}
			fmt.Println("Error: Wrong date format")
			returnJs(w, resText)
		}

		if PNOdate != today {
			PNOyear  := PNO[2:6]
			PNOmonth := PNO[6:8]
			PNOday   := PNO[8:10]
			fmt.Println("PNOdate != today")
			insertDate := PNOyear + "-" + PNOmonth + "-" + PNOday + " " + tdTime
			fmt.Println("InsertDate:", insertDate)
			insertQuery := map[string]interface{}{
				"PNO":PNO,
				"model":model,
				"sales_mgr": salesMgr,
				"qty":qty,
				"customer":customer,
				"location":location,
				"status":status,
				"created_at": insertDate,
			}
			for kname, val := range insertQuery {
				if val == "" {
					resText := map[string]string {
						"err": "Empty form data",
					}
					fmt.Println("Error: Empty form data:", kname)
					returnJs(w, resText)
				}
			}
			insertFeedback := mysql.InsertInto(tbname["picklist"], insertQuery).Use(db)
			if err != nil {
				fmt.Println("Commit error:", err)
			}
			returnJson(w, insertFeedback)
		} else if PNOdate == today {
			fmt.Println("PNOdate == today")
			
			insertQuery := map[string]interface{}{
				"PNO":PNO,
				"model":model,
				"sales_mgr": salesMgr,
				"qty":qty,
				"customer":customer,
				"location":location,
				"status":status,
			}

			for kname, val := range insertQuery {
				if val == "" {
					resText := map[string]string {
						"err": "Empty form data",
					}
					fmt.Println("Error: Empty form data:", kname)
					returnJs(w, resText)
				}
			}

			insertFeedback := mysql.InsertInto(tbname["picklist"], insertQuery).Use(db)
			if err != nil {
				fmt.Println("Commit error:", err)
			}
			returnJson(w, insertFeedback)
		}
		
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
		PNO      := r.FormValue("PNO")
		model    := r.FormValue("model")
		qty      := r.FormValue("qty")
		customer := r.FormValue("customer")
		location := r.FormValue("location")
		status   := r.FormValue("status")
		salesMgr := r.FormValue("sales_mgr")
		
		tx, ctx := mysql.Begin(db)
		updateInfo := map[string]interface{} {
			"PNO":PNO,
			"model":model,
			"sales_mgr": salesMgr,
			"qty":qty,
			"customer":customer,
			"location":location,
			"status":status,
		}
		mysql.Update(tbname["picklist"], false).Set(updateInfo).Where("PID",PID).With(tx, ctx)
		fmt.Printf("[%-18s] Getting PID=%s from UPDATEd\n", "PickListUpdate", PID);
		updatedPicked := mysql.
			Select("PNO", "model","sales_mgr", "qty", "customer", "location", "status").
			From(tbname["picklist"]).
			Where("PID", PID).
			With(tx, ctx)
		dbCommits["PickList"] = tx
		fmt.Println("updatedPicked:", updatedPicked)
		returnJson(w, updatedPicked)
	}
}
