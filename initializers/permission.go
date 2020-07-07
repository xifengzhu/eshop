package initializers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var Permissions []string

func init() {
	jsonFile, err := os.Open("./conf/permissions.json")

	if err != nil {
		fmt.Println("load permission json error: ", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Permissions)
}
