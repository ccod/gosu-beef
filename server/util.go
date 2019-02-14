package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondJSON(w http.ResponseWriter, r *http.Request, d interface{}) {
	w.Header().Set("Content-type", "application/json")

	j, err := json.Marshal(d)
	if err != nil {
		fmt.Printf("Failed to generate JSON: %s\n", err)
		return
	}

	w.Write(j)
}
