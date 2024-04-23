package handler

import (
	"encoding/json"
	"github.com/mr-Evgeny/go_final_project/database"
	"github.com/mr-Evgeny/go_final_project/model"
	"net/http"
	"time"
)

const RowsLimit int = 50

func TaskList(w http.ResponseWriter, r *http.Request) {
	tasks := []model.ToDo{}

	query := r.URL.Query()
	search := query.Get("search")
	is_searched := false
	//checking if search is date
	if len(search) == 10 && search[2] == '.' && search[5] == '.' {
		search_date, err := time.Parse("02.01.2006", search)
		if err == nil {
			is_searched = true
			resultDB := database.DB.Db.Limit(RowsLimit).Order("date, id").Where("date = ?", search_date.Format("20060102")).Find(&tasks)
			if resultDB.Error != nil {
				ErrorJson(w, resultDB.Error.Error())
				return
			}
		}
	}
	//text search
	if !is_searched && len(search) > 0 {
		//PRAGMA case_sensitive_like=OFF;
		is_searched = true
		resultDB := database.DB.Db.Limit(RowsLimit).Order("date, id").Where("title LIKE ?", "%"+search+"%").Or("comment LIKE ?", "%"+search+"%").Find(&tasks)
		if resultDB.Error != nil {
			ErrorJson(w, resultDB.Error.Error())
			return
		}
	}
	//list of tasks
	if !is_searched {
		resultDB := database.DB.Db.Limit(RowsLimit).Order("date, id").Find(&tasks)
		if resultDB.Error != nil {
			ErrorJson(w, resultDB.Error.Error())
			return
		}
	}

	resp, err := json.Marshal(struct {
		T []model.ToDo `json:"tasks"`
	}{tasks})
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
