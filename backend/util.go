package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getHostAndPortBak(addr string, port int) string {
	return fmt.Sprintf("%v:%v", addr, port)
}

func HTTP404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func HTTP200(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "Status OK"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}

func HTTP500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
