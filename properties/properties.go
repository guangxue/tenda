package properties

import (
	"encoding/json"
	"fmt"
	"os"
)

var propsFile string = "properties.json"
var properties = make(map[string]string)

func GetProperty(propName string) string {
	propBlob, err := os.ReadFile(propsFile)
	if err != nil {
		fmt.Println("Read Properties Error: ", err)
	}

	// Read JSON and store in `properties`
	err = json.Unmarshal(propBlob, &properties)

	if err != nil {
		fmt.Println(" JSON Unmarshal Error: ", err)
	}

	if len(propName) != 0 {
		if theProp, found := properties[propName]; found {
			return theProp
		} else {
			return ""
		}
	} else {
		return ""
	}

}
