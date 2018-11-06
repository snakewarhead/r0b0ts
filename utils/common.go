package utils

import (
	"encoding/json"
	"time"
)

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

func MustNotEmpty(args ...string) bool {
	for _, s := range args {
		if len(s) == 0 {
			return false
		}
	}
	return true
}

func ToJsonBelievably(v interface{}) string {
	s, err := json.Marshal(v)
	if err != nil {
		Logger.Error("parse transaction to json error --- %v", err)
		return ""
	}
	return string(s)
}
