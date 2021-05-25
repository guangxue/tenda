package tenda

import (
	"fmt"
	"net/http"
	"strings"
	"github.com/guangxue/webapps/mysql"
)

func Model(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("[%-18s] search model path:%s\n", "Model", r.URL.Path)
	searchLocation := r.URL.Query().Get("location");
	fmt.Printf("[%-18s] searchLocation:%s\n", "Model", searchLocation)
	searchModel := strings.TrimPrefix(r.URL.Path, "/tenda/api/model/")
	fmt.Printf("[%-18s] search model:%s\n", "Model", searchModel)
    if strings.Contains(searchModel, "/") {
        searchModel = ""
    }

	if searchModel != "" && searchLocation != "" {
		allModels := mysql.
			Select("sid","model","location","unit","cartons","boxes","total").
			From(tbname["stock_updated"]).
			Where("model", searchModel).
			AndWhere("location", "=", searchLocation).
		Use(db)
		returnJson(w, allModels)
	}

	if searchModel == "" && searchLocation == "" {
		modelNames := mysql.
			SelectDistinct("model").
			From(tbname["stock_updated"]).
		Use(db);
		returnJson(w, modelNames)
	}

	if searchModel != "" && searchLocation == "" {
		allModels := mysql.Select("model", "location", "unit", "cartons", "boxes", "total").
			From(tbname["stock_updated"]).
			Where("model", searchModel).
			Use(db)
	    returnJson(w, allModels)
	}

	if searchModel == "" && searchLocation != "" {
		allModels := mysql.
			Select("model","location","unit","cartons","boxes","total").
			From(tbname["stock_updated"]).
			Where("location", searchLocation).
		Use(db)
		for _, val := range allModels {
			fmt.Println(val)
		}
		returnJson(w, allModels)
	}
}
