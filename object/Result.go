package object

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Code    int
	Count   int
	Message string
	Data    interface{}
}

func (this Result) SetData(d interface{}) Result {
	if d != nil {
		this.Data = d
	}
	return this
}

func (this Result) SetCount(c int) Result {
	this.Count = c
	return this
}

func (this Result) String() string {
	b, err := json.Marshal(this)
	if err != nil {
		fmt.Println("error:", err)
	}
	return string(b)
}

func ResultNew(flag bool, success, fail string) Result {
	if flag {
		return Result{Code: 1, Message: success}
	} else {
		return Result{Code: -1, Message: success}
	}
}

func ResultNewSuccess(msg string, d interface{}) Result {
	return Result{Code: 1, Message: msg, Data: d}
}

func ResultNewFail(msg string, d interface{}) Result {
	return Result{Code: -1, Message: msg, Data: d}
}
