package tenda

import (
	"fmt"
	"net/http"
	"github.com/guangxue/webapps/mysql"
)

func Models(w http.ResponseWriter, r *http.Request) {

	queryModel    := r.URL.Query().Get("model");
	queryLocation := r.URL.Query().Get("location");

	fmt.Printf("[%-18s] query Model:%s\n", "Models", queryModel)
	fmt.Printf("[%-18s] query Location:%s\n", "Models", queryLocation)


    // get all models
	if len(queryModel) == 0 && len(queryLocation) == 0{
        modelNames := mysql.
        	SelectDistinct("model").
        	From(tbname["stock_updated"]).Use(db);
        returnJson(w, modelNames)
	}

    // get one model
	if len(queryModel) > 0 {
        allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").From(tbname["stock_updated"]).Where("model", queryModel).Use(db)
	    returnJson(w, allModels)
	}

	// get location data
	if len(queryLocation) > 0 {
		allModels := mysql.
			Select("model", "location", "unit", "cartons", "boxes", "total").
			From(tbname["stock_updated"]).
			Where("location", queryLocation).
		Use(db)
		returnJson(w, allModels)
	}
}