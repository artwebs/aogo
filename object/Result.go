package object

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	data map[string]interface{}
}

func (this Result) Append(key string, val interface{}) Result {
	this.data[key] = val
	return this
}

func (this Result) String() string {
	b, err := json.Marshal(this.data)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}

func ResultNew(flag bool, success, fail string) Result {
	rs := Result{data: map[string]interface{}{}}
	if flag {
		rs.data["code"] = 1
		rs.data["message"] = success
	} else {
		rs.data["code"] = 1
		rs.data["message"] = success
	}
	return rs
}

func ResultNewSuccess(msg string) Result {
	return Result{data: map[string]interface{}{"code": 1, "message": msg}}
}

func ResultNewFail(msg string) Result {
	return Result{data: map[string]interface{}{"code": 0, "message": msg}}
}
