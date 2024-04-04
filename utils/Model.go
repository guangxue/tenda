package tenda

import (
	"fmt"
	"net/http"
	"strings"
	"tenda/mysql"
)

func Model(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%-18s] request path:%s\n", " -- Model.go", r.URL.Path)
	searchLocation := r.URL.Query().Get("location");
	fmt.Printf("[%-18s] ?location = %s\n", " -- Model.go", searchLocation)
	searchModel := strings.TrimPrefix(r.URL.Path, "/api/model/")
	
  if strings.Contains(searchModel, "/") {
    searchModel = ""
  }
  fmt.Printf("[%-18s] /model = %s\n", " -- Model.go", searchModel)
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
		allModels := mysql.Select("SID","model", "location", "unit", "cartons", "boxes", "total").
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
