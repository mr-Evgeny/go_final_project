package handler

import (
	"fmt"
	"net/http"
	"time"
)

func NextDate(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	now_str := query.Get("now")
	if len(now_str) != 8 {
		http.Error(w, "field 'now' not specified", http.StatusBadRequest)
		return
	}
	now, err := time.Parse("20060102", now_str)
	if err != nil {
		http.Error(w, "field 'now' has wrong format", http.StatusBadRequest)
		return
	}
	date := query.Get("date")
	if len(date) != 8 {
		http.Error(w, "field 'date' not specified", http.StatusBadRequest)
		return
	}
	repeat := query.Get("repeat")

	//--> datecalc.go
	nextDate, err := getNextDate(now, date, repeat)
	if err != nil {
		http.Error(w, "bad date: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, nextDate)
}
