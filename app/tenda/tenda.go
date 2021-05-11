package tenda

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	_ "strings"
	"encoding/json"
	"time"
	"strconv"
	"html/template"
	"github.com/guangxue/webapps/mysql"
)

type pickcolumns struct {
	PID      int
	Model    string
	Qty      int
	Location string
}
type updateModel struct {
	Location string
	Model    string
	Unit     int
	Cartons  int
	Boxes    int
	Total    int
}

var db = mysql.Connect("tenda");

func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func WriteJSON(w http.ResponseWriter, returnRows []map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	// respJSON, err := json.Marshal(returnRows)
	// if err != nil {
	// 	fmt.Println("returnRows JSON Marshal error: ", err)
	// }
	// w.Write(respJSON)
	json.NewEncoder(w).Encode(returnRows)
}

func TimeNow() string {
	ts := time.Now()
	now := ts.Format("2006-01-02 15:04:05")
	return now
}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	tmplpath := "templates/tenda/" + templateName
	tmpl, err := template.ParseFiles(tmplpath, "templates/tenda/base.html","templates/tenda/nav.html")
	if err != nil {
		fmt.Println("template parsing errors: ", err)
	}
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println("template executing errors: ", err)
	}
}

func RenderHandler(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// c, err := r.Cookie("gossessid")
		// if err != nil {
		// 	http.Redirect(w, r, "/tenda", http.StatusSeeOther)
		// 	return
		// }
		// fmt.Println("c.Value:", c.Value)
		tmplpath := "templates/tenda/" + templateName
		tmpl, err := template.ParseFiles("templates/tenda/base.html", "templates/tenda/nav.html", tmplpath)
		if err != nil {
			fmt.Println("template parsing errors: ", err)
		}
		err = tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
			fmt.Println("template executing errors: ", err)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	render(w, "login.html", nil)
}
func Models(w http.ResponseWriter, r *http.Request) {

    // Set Header for json HTTP response
	w.Header().Set("Content-Type", "application/json")

	queryModel    := r.URL.Query().Get("model");
	queryLocation := r.URL.Query().Get("location");

	fmt.Printf("[%-18s] Request Path:%s\n", "Models", r.URL.Path)
	fmt.Printf("[%-18s] query Model:%s\n", "Models", queryModel)
	fmt.Printf("[%-18s] query Location:%s\n", "Models", queryLocation)


    // get all models
	if len(queryModel) == 0 && len(queryLocation) == 0{
        modelNames := mysql.SelectDistinct("model").From("stock_updated").Use(db);
	    ModelNamesJSON, err := json.Marshal(modelNames)
	    if err != nil {
	    	fmt.Println("[Models] ModelsJson error: ", err)
	    }
		w.Write(ModelNamesJSON)
	}

    // get one model
	if len(queryModel) > 0 {
        allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From("stock_updated").Where("model", queryModel).Use(db)
	    ModelsJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
        w.Write(ModelsJSON)
	}

	// get location data
	if len(queryLocation) > 0 {
		allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From("stock_updated").Where("location", queryLocation).Use(db)
		LocationJSON, err := json.Marshal(allModels)
	    ErrorCheck(err)
	    fmt.Printf("[%-18s] LocationJSON:%v\n", "Models", string(LocationJSON))
 	    w.Write(LocationJSON)
	}
}

func Locations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryModel := r.URL.Query().Get("model");
	
	if len(queryModel) > 0 {
		allLocations := mysql.Select("location").From("stock_updated").Where("model", queryModel).Use(db)
		LocationJSON, err := json.Marshal(allLocations)
	    ErrorCheck(err)
	    for _, val := range allLocations {
	    	fmt.Printf("[%-18s] %v\n","Locations", val)
	    }
 	    w.Write(LocationJSON)
	}
}

func PickList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		date := r.URL.Query().Get("date");
		status := r.URL.Query().Get("status");
		PID := r.URL.Query().Get("PID");

		fmt.Printf("[%-18s] date  :%s\n", "PickList", date);
		fmt.Printf("[%-18s] status:%s\n", "PickList", status);
		fmt.Printf("[%-18s] PID   :%s\n", "PickList", PID);

		if status == "completed_at" {
			// Weekly completed orders
			startDate := fmt.Sprintf("'%s'", date)
			endDate := fmt.Sprintf("date_add('%s', INTERVAL 7 DAY)", date)
			allPicked := mysql.Select("LID", "location", "model", "unit", "cartons", "boxes", "total", "completed_at").From("last_updated").WhereBetween("completed_at", startDate, endDate).Use(db)
			json.NewEncoder(w).Encode(allPicked)
		} else if len(PID) > 0 && len(status) > 0 {
			// Picked orders
			allPicked := mysql.Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").From("picklist").Where("PID",PID).AndWhere("status", "=",status).Use(db)
			json.NewEncoder(w).Encode(allPicked)
		} else if len(date) > 0 && status == "created_at" { 
			// Weekly picked
			stmt := fmt.Sprintf("WITH weeklypicked AS (select model, qty, customer, location, status, created_at FROM picklist Where created_at BETWEEN '%s' AND date_add('%s', INTERVAL 7 DAY)) SELECT model, SUM(qty) as total from weeklypicked group by model", date, date)
			allPicked := mysql.SelectRaw(stmt, "model", "total").Use(db)
			fmt.Printf("[%-18s] PID   :%s %v\n", "weeklypicked:allpicked:", allPicked)
			json.NewEncoder(w).Encode(allPicked)
		} else {
			odate := ""
			if date == "" {
				odate = "%"
			} else {
				odate = date + "%"
			}
			allPicked := mysql.Select("PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "updated_at").From("picklist").WhereLike("created_at",odate).AndWhere("status", "=",status).Use(db)
			json.NewEncoder(w).Encode(allPicked)
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
		insertFeedback := mysql.InsertInto("picklist", insertQuery).Use(db)
		
		respJSON, err := json.Marshal(insertFeedback)
	    if err != nil {
	    	fmt.Println("respJSON error: ", err)
	    }
		w.Write(respJSON)
	}
}

func CompletePickList (w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			fmt.Printf("[%-18s] Form parse error:\n", "CompletePickList", err)
		}
		/* {pickDate}   - POST request from Ajax
		/* {pickStatus} - POST request from Ajax
		/* 1. Parse POST form data {pickDate}, {pickStatus} */
		pickDate := r.FormValue("pickDate")
		pickStatus := r.FormValue("pickStatus")
		lastSaturday := r.FormValue("lastSaturday")
		
		fmt.Printf("[%-18s] Pick date  :%v\n", "CompletePickList",pickDate)
		fmt.Printf("[%-18s] pick status:%v\n", "CompletePickList",pickStatus)
		fmt.Printf("[%-18s] pick lastSaturday:%v\n", "CompletePickList",lastSaturday)

		p := pickcolumns{}
		pcols := []pickcolumns{}

		/* if {pickDate} is empty, no db action needed. */
		if pickStatus == "Updated"  || pickDate == "" {
			fmt.Printf("[%-18s] *Stmt Err*: Invalid update statement for complete orders, want {pickStatus} and {pickDate}.", "CompletePickList")
			fmt.Printf("[%-18s] *Stmt Err*: Want {pickStatus} and {pickDate} to complete.", "CompletePickList")
			return 
		}

		/* 2. SELECT FROM `picklist` according {pickDate} and {pickStatus} which from POST data. */
		/*   {p}   - Scaned single row */
		/* {pcols} - array of {p}  */
		sqlstmt := "SELECT PID, model, qty, location FROM picklist WHERE created_at LIKE '"+pickDate+"%' AND status ='"+pickStatus+"'"
		rows, err := db.Query(sqlstmt)
		if err != nil {
			fmt.Printf("[%-18s] *SELECT Err*:%v\n", "CompletePickList", err)
		}

		for rows.Next() {
			err := rows.Scan(&p.PID, &p.Model, &p.Qty, &p.Location)
			if err != nil {
				fmt.Printf("[%-18s]: dbColumns Scan error:%v\n", "CompletePickList",err)
			}
			pcols = append(pcols, p)
		}
		/* 3. Start calculate {cartons}, {boxes}, {total} to UPDATE  */
		upModels := []updateModel{}

		for _, p := range pcols {
			fmt.Printf("[%-18s] Model    :%s\n", "CompletePickList *",p.Model)
			fmt.Printf("[%-18s] Location :%s\n", "CompletePickList",p.Location)
			
			unit := 0
			oldCartons := 0
			oldBoxes := 0
			oldTotals := 0
			/* 4.0 Get original {total}, {unit} from `stock_updated` */
			err := db.QueryRow("SELECT unit, cartons, boxes, total FROM stock_updated WHERE model=? AND location=?", p.Model, p.Location).Scan(&unit, &oldCartons, &oldBoxes, &oldTotals)
			switch {
				case err == sql.ErrNoRows:
					fmt.Printf("[db *ERR*] no `oldTotals` for model %s\n", p.Model)
				case err != nil:
					fmt.Printf("[db *ERR*] query error: %v\n", err)
				default:
					fmt.Printf("[%-18s] oldTotals are %d\n", "CompletePickList",oldTotals)
			}
			// {p.Qty}: quantity picked
			fmt.Printf("[%-18s] pick qty:%d\n", "CompletePickList",p.Qty)

			newTotal := oldTotals - p.Qty
			fmt.Printf("[%-18s] *NEW Total:%d\n", "CompletePickList",newTotal)
			fmt.Printf("[%-18s] *unit are %d\n", "CompletePickList",unit)

			/* 4.1 if unit = 0, then {newCartons} = {newBox} */
			newCartons := 0
			newBoxes   := newTotal;

			/* 4.2 {newTotal}  : {total} - {p.Qty} */
			/* 4.3 {newCartons}: {newCartons}/{unit} */
			/* 4.4 {newBoxes}  : ({newCartons}/{unit} - {newCartons} )*{unit} */
			if unit > 1 {
				newCartons = newTotal/unit
				fmt.Printf("[%-18s] *NEW Cartons:%d\n", "CompletePickList",newCartons)
				newBoxesFrac := float64(newTotal)/float64(unit) - float64(newCartons)
				fmt.Printf("[%-18s] *NEW BoxeFrac:%f\n", "CompletePickList",newBoxesFrac)
				newBoxesFrac = newBoxesFrac * float64(unit)
				newBoxes = int(math.Round(newBoxesFrac))
				fmt.Printf("[%-18s] *NEW Boxes:%d\n", "CompletePickList",newBoxes)
			} 
			
			upModel := updateModel{p.Location, p.Model, unit, newCartons, newBoxes, newTotal}
			upModels = append(upModels, upModel)
			// ----------------------------------------------------
			/* 5. Update `stock_update` table first */
			updateStockUpdate := map[string]interface{} {
				"cartons": newCartons,
				"boxes"  : newBoxes,
				"total"  : newTotal,
			}
			mysql.Update("stock_updated",false).Set(updateStockUpdate).Where("model", p.Model).AndWhere("location", "=",p.Location).Use(db)
			// -----------------------------------------------------

			/* 6. Update `picklist` table status to 'Complete' */
			updatePLInfo := map[string]interface{} {"status":"Complete"}
			mysql.Update("picklist", false).Set(updatePLInfo).Where("PID", strconv.Itoa(p.PID)).Use(db)

			/* 7. Check if model (that is completed) already exists in the table `last_updated` */
			/* 7.1            IF EXISTS, UPDATE it,
			 * ....otherwise, INSERT: new data  */
			ckmodel := ""
			Checkerr := db.QueryRow("SELECT model FROM last_updated WHERE model=? AND location=? AND completed_at > ?", p.Model, p.Location, lastSaturday).Scan(&ckmodel)
			switch {
			case Checkerr == sql.ErrNoRows:
				fmt.Printf("[%-18s] NO ROWS return from last_updated for `%s`, then INSERT\n", "*db Rows*", p.Model)
				insertValues := map[string]interface{} {
					"location": p.Location,
					"model"   : p.Model,
					"unit"    : unit,
					"cartons" : newCartons,
					"boxes"   : newBoxes,
					"total"   : newTotal,
					"completed_at": TimeNow(),
				}
				mysql.InsertInto("last_updated",insertValues).Use(db);
			case Checkerr != nil:
				fmt.Printf("[db *ERR*] query error: %v\n", err)
			case len(ckmodel) > 0:
				fmt.Printf("[%-18s] FOUND `model`:%s , then UPDATE\n", "*db Rows*", ckmodel)
				updateValues := map[string]interface{} {
					"location": p.Location,
					"model"   : p.Model,
					"unit"    : unit,
					"cartons" : newCartons,
					"boxes"   : newBoxes,
					"total"   : newTotal,
					"completed_at": TimeNow(),
				}
				mysql.Update("last_updated", false).Set(updateValues).Where("model", p.Model).AndWhere("location", "=",p.Location).AndWhere("completed_at", ">", lastSaturday).Use(db)
			}
		}
		// end for loop
	}
}

func Stocktakes (w http.ResponseWriter, r *http.Request) {
	tbName := r.URL.Query().Get("tbname")

	if len(tbName) > 0 {
		allstocks := mysql.Select("SID", "location", "model", "unit", "cartons", "boxes","total", "kind", "notes").From(tbName).Use(db)
		json.NewEncoder(w).Encode(allstocks)
	} else {
		render(w, "stocktakes.html", nil)
	}
	
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	SID := r.URL.Query().Get("SID")
	tbName := r.URL.Query().Get("tbname")

	if SID != "" && r.Method == http.MethodGet {
		currentStockToUpdate := mysql.Select("SID", "location", "model", "unit", "cartons", "boxes","total").From(tbName).Where("SID", SID).Use(db);
		fmt.Println("currentStockToUpdate:", currentStockToUpdate)
		render(w, "update_stock.html", currentStockToUpdate[0])
	}
}

func UpdatePickList(w http.ResponseWriter, r *http.Request) {

	dbPickedInfo := map[string]string{}

	if r.Method == "GET" {
		queryPID := r.URL.Query().Get("PID");
		fmt.Printf("[%-18s] PID:%s\n", "UpdatePickList",queryPID)
		currentPID := mysql.Select("PNO", "model", "qty", "customer", "location", "status").From("picklist").Where("PID", queryPID).Use(db)
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
		render(w, "update_picked.html", dbPickedInfo)
	}

	if r.Method == "POST" {
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
		status := r.FormValue("statusoption")
		fmt.Printf("PNO:%s\nmodel:%s\nqty:%s\ncustomer:%s\nlocation:%s\nstatus:%s",PNO, model, qty, customer,location,status)
		fmt.Println("dbPickedInfo::", dbPickedInfo)
		
		updateInfo := map[string]interface{} {
			"PNO":PNO,
			"model":model,
			"qty":qty,
			"customer":customer,
			"location":location,
			"status":status,
		}
		rowsAffaced := mysql.Update("picklist", false).Set(updateInfo).Where("PID",PID).Use(db)
		json.NewEncoder(w).Encode(rowsAffaced)
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
		mysql.DeleteFrom("picklist", false).Where("PID", PID).Use(db)
	}

	if status == "Complete" {
		// rollback
	}
}

func AddStock(w http.ResponseWriter, r *http.Request) {
	render(w, "addstock.html", nil)
}