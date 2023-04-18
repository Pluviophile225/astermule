package parser

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Pluviophile225/astermule/tools/whiteboard"
)

func Test_1(t *testing.T) {
	var param map[string]string
	err := json.Unmarshal([]byte(`{"id":"S(\"init\", \"userId\")"}`), &param)
	if err != nil {
		fmt.Println("err:", err.Error())
	}
	fmt.Println(param)
}
func Test_2(t *testing.T) {
	param := make(map[string]string)
	var ParamFormat = "{\"id\":\"S(\\\"init\\\",\\\"userId\\\")\"}"
	var ActionMaps = map[string]map[string]interface{}{
		"init": {
			"userId": "123",
		},
	}
	json.Unmarshal([]byte(ParamFormat), &param)
	fmt.Println("ActionMap: ", ActionMaps)
	fmt.Println("node param format: ", ParamFormat)
	fmt.Println("rule param:", param)
	data, err := whiteboard.Bend(param, ActionMaps, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("node result data:", data)
}
