package main

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func handleDefault(w http.ResponseWriter, httpStatus int, status bool, msg string) {
	rs, err := jsoniter.Marshal(JsonRs{httpStatus, status, msg})
	if err != nil {
		panic("handleDefault: - " + err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(rs)
}
