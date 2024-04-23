package handler

import (
	"encoding/json"
	"net/http"
)

func ErrorJson(w http.ResponseWriter, msg string) {
	resp, _ := json.Marshal(struct {
		E string `json:"error"`
	}{msg})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(resp)
}

func Api(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/tasks" && r.Method == "GET" {
		TaskList(w, r)
		return
	}
	if r.URL.Path == "/api/task/done" && r.Method == "POST" {
		TaskDone(w, r)
		return
	}
	if r.URL.Path == "/api/task" {
		switch r.Method {
		case "POST":
			TaskAdd(w, r)
			return
		case "GET":
			TaskInfo(w, r)
			return
		case "PUT":
			TaskEdit(w, r)
			return
		case "DELETE":
			TaskDelete(w, r)
			return
		}
	}
	http.Error(w, "Unknown target", http.StatusBadRequest)
	return
}
