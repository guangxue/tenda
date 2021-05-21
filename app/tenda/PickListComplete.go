package tenda

import (
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/guangxue/webapps/mysql"
)

func PickListComplete (w http.ResponseWriter, r *http.Request) {
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

		fmt.Printf("[%-18s] Pick date   :%v\n", "CompletePickList",pickDate)
		fmt.Printf("[%-18s] pick status :%v\n", "CompletePickList",pickStatus)
		fmt.Printf("[%-18s] pick lastSaturday :%v\n", "CompletePickList",lastSaturday)


		p := pickcolumns{}
		pcols := []pickcolumns{}

		/* if {pickDate} is empty, no db action needed. */
		// if pickStatus == "Updated"  || pickDate == "" {
			
		// 	fmt.Printf("[%-18s] *Stmt Err*: Invalid update statement for complete orders, want {pickStatus} and {pickDate}.\n", "CompletePickList")
		// 	fmt.Printf("[%-18s] *Stmt Err*: Want {pickStatus} and {pickDate} to complete.\n", "CompletePickList")
		// 	stmtErr := map[string]string{
		// 		"stmtErr": "Need {pickStatus} and {pickDate} to complete picklist",
		// 	}
		// 	stmtErrors := []map[string]string{}
		// 	stmtErrors = append(stmtErrors, stmtErr)
		// 	returnJson(w, stmtErrors)
		// 	return
		// }

		/* 2. SELECT FROM `picklist` according {pickDate} and {pickStatus} which from POST data. */
		/*   {p}   - Scaned single row */
		/* {pcols} - array of {p}  */

		stmt := fmt.Sprintf("SELECT PID, model, qty, location FROM %s WHERE created_at LIKE %q AND status =%q", tbname["picklist"], pickDate+"%", pickStatus)
		fmt.Printf("[%-18s] Select Stmt :%s\n", "CompletePickList", stmt)
		// sqlstmt := "SELECT PID, model, qty, location FROM picklist WHERE created_at LIKE '"+pickDate+"%' AND status ='"+pickStatus+"'"
		/************************/
		//						//
		tx, ctx := mysql.Begin(db)
		//						//
		/************************/
		rows, err := tx.QueryContext(ctx, stmt)
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
		completeInfos := []map[string]string{}
		upModels := []updateModel{}
		for i, p := range pcols {
			fmt.Printf("[%-18s] index    :%v\n", "CompletePickList *",i)
			fmt.Printf("[%-18s] Model    :%s\n", "CompletePickList *",p.Model)
			fmt.Printf("[%-18s] Location :%s\n", "CompletePickList",p.Location)
			completeInfo := map[string]string{}
			completeInfo["model"] = p.Model
			completeInfo["location"] = p.Location
			unit := 0
			oldCartons := 0
			oldBoxes := 0
			oldTotals := 0
			/* 4.0 Get original {total}, {unit} from `stock_updated` */
			stmt := fmt.Sprintf("SELECT unit, cartons, boxes, total FROM %s WHERE model=? AND location=?", tbname["stock_updated"])
			err := tx.QueryRowContext(ctx,stmt, p.Model, p.Location).Scan(&unit, &oldCartons, &oldBoxes, &oldTotals)
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
			completeInfo["pickQty"]= strconv.Itoa(p.Qty)

			newTotal := oldTotals - p.Qty
			fmt.Printf("[%-18s] *NEW Total:%d\n", "CompletePickList",newTotal)
			fmt.Printf("[%-18s] *unit are %d\n", "CompletePickList",unit)
			completeInfo["oldTotal"]= strconv.Itoa(oldTotals)
			completeInfo["unit"] = strconv.Itoa(unit)

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
			completeInfo["newCartons"] = strconv.Itoa(newCartons)
			completeInfo["newBoxes"] = strconv.Itoa(newBoxes)
			completeInfo["newTotal"] = strconv.Itoa(newTotal)
			updateStockUpdate := map[string]interface{} {
				"cartons": newCartons,
				"boxes"  : newBoxes,
				"total"  : newTotal,
			}
			mysql.
				Update(tbname["stock_updated"],false).
				Set(updateStockUpdate).
				Where("model", p.Model).
				AndWhere("location", "=",p.Location).
			With(tx,ctx)
			// -----------------------------------------------------

			/* 6. Update `picklist` table status to 'Complete' */
			pickListSet := map[string]interface{} {"status":"Complete"}
			mysql.
				Update(tbname["picklist"], false).
				Set(pickListSet).
				Where("PID", strconv.Itoa(p.PID)).
			With(tx,ctx)

			/* 7. Check if model (that is completed) already exists in the table `last_updated` */
			/* 7.1            IF EXISTS, UPDATE it,
			 * ....otherwise, INSERT: new data  */
			
			existModelId := 0
			stmt = fmt.Sprintf("SELECT LID FROM %s WHERE model=? AND location=? AND completed_at > ?", tbname["last_updated"])
			Checkerr := tx.QueryRowContext(ctx,stmt, p.Model, p.Location, lastSaturday).Scan(&existModelId)
			switch {
			case Checkerr == sql.ErrNoRows:
				fmt.Printf("[%-18s] NO ROWS return from last_updated for `%s`, then INSERT\n", "*db Rows*", p.Model)
				completeInfo["sqlinfo"] = "INSERT into last_updated"
				insertValues := map[string]interface{} {
					"location": p.Location,
					"model"   : p.Model,
					"unit"    : unit,
					"cartons" : newCartons,
					"boxes"   : newBoxes,
					"total"   : newTotal,
					"completed_at": TimeNow(),
				}
				mysql.InsertInto(tbname["last_updated"],insertValues).With(tx, ctx)
			case Checkerr != nil:
				fmt.Printf("[db *ERR*] query error: %v\n", err)
			case  existModelId > 0:
				fmt.Printf("[%-18s] Exits `model id`:%s , then UPDATE\n", "*db Rows*", existModelId)
				completeInfo["sqlinfo"] = "UPDATE last_update"
				updateValues := map[string]interface{} {
					"location": p.Location,
					"model"   : p.Model,
					"unit"    : unit,
					"cartons" : newCartons,
					"boxes"   : newBoxes,
					"total"   : newTotal,
					"completed_at": TimeNow(),
				}
				mysql.
					Update(tbname["last_updated"], false).
					Set(updateValues).
					Where("model", p.Model).
					AndWhere("location", "=",p.Location).
					AndWhere("completed_at", ">", lastSaturday).
				With(tx, ctx)
			}
			completeInfos = append(completeInfos, completeInfo)
		}
		dbCommits["CompletePickList"] = tx
		json.NewEncoder(w).Encode(completeInfos)
	}
}
