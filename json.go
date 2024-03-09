package main

import (
	"github.com/gogf/gf/v2/encoding/gjson"
)

func ToJsonString(target interface{}) string {
	if target == nil {
		return ""
	}
	marshal, err := gjson.Marshal(target)
	if err != nil {
		return ""
	}
	return string(marshal)
}
