package mysql

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
 )

type DbInfo struct {
  Usr string `json:"usr"`
  Pwd string `json:"pwd"`
}


func dbAuthStr() string {

  // reading config file
  configData, err := ioutil.ReadFile("./dsn.json") 
  if err != nil {
    fmt.Print("[Config File Read Error] ", err)
  }

  // json data
  var db DbInfo

  err = json.Unmarshal(configData, &db)

  if err != nil {
    fmt.Println("Error [json]", err)
  }

  return db.Usr+":"+db.Pwd+"@/"
}

