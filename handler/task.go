package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/mr-Evgeny/go_final_project/database"
	"github.com/mr-Evgeny/go_final_project/model"
	"net/http"
	"strconv"
	"time"
)

func taskSet(r *http.Request, task *model.ToDo) error {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return err
	}
	if len(task.Title) == 0 {
		return errors.New("No title specified")
	}
	if len(task.Date) > 0 {
		date_tochk, err := time.Parse("20060102", task.Date)
		if err != nil {
			return err
		}
		now := Today()
		if date_tochk.Before(now) {
			if len(task.Repeat) > 0 {
				nextDate, err_calc := getNextDate(now, task.Date, task.Repeat)
				if err_calc != nil {
					return err_calc
				}
				task.Date = nextDate
			} else {
				task.Date = now.Format("20060102")
			}
		}
	} else {
		task.Date = Today().Format("20060102")
	}
	return nil
}

func TaskAdd(w http.ResponseWriter, r *http.Request) {
	task := new(model.ToDo)
	err := taskSet(r, task)
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	resultDB := database.DB.Db.Create(&task)
	if resultDB.Error != nil {
		ErrorJson(w, resultDB.Error.Error())
		return
	}

	json_resp, err := json.Marshal(struct {
		Id uint `json:"id"`
	}{task.ID})
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json_resp)
}

func TaskEdit(w http.ResponseWriter, r *http.Request) {
	task := new(model.ToDo)
	err := taskSet(r, task)
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}
	if task.ID == 0 {
		ErrorJson(w, "task id not specified")
		return
	}
	resultDB := database.DB.Db.Save(&task)
	if resultDB.Error != nil {
		ErrorJson(w, resultDB.Error.Error())
		return
	}
	json_resp, _ := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}

func TaskInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if len(query.Get("id")) == 0 {
		ErrorJson(w, "ID not specified")
		return
	}
	id, err := strconv.Atoi(query.Get("id"))
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	task := new(model.ToDo)

	resultDB := database.DB.Db.First(&task, id)
	if resultDB.Error != nil {
		ErrorJson(w, resultDB.Error.Error())
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if len(query.Get("id")) == 0 {
		ErrorJson(w, "ID not specified")
		return
	}
	id, err := strconv.Atoi(query.Get("id"))
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	resultDB := database.DB.Db.Delete(&model.ToDo{}, id)
	if resultDB.Error != nil {
		ErrorJson(w, resultDB.Error.Error())
		return
	}

	json_resp, _ := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}

func TaskDone(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if len(query.Get("id")) == 0 {
		ErrorJson(w, "ID not specified")
		return
	}
	id, err := strconv.Atoi(query.Get("id"))
	if err != nil {
		ErrorJson(w, err.Error())
		return
	}

	task := new(model.ToDo)

	resultDB := database.DB.Db.First(&task, id)
	if resultDB.Error != nil {
		ErrorJson(w, resultDB.Error.Error())
		return
	}
	if len(task.Repeat) == 0 {
		resultDB := database.DB.Db.Delete(&model.ToDo{}, id)
		if resultDB.Error != nil {
			ErrorJson(w, resultDB.Error.Error())
			return
		}
	} else {
		nextDate, err_calc := getNextDate(time.Now(), task.Date, task.Repeat)
		if err_calc != nil {
			ErrorJson(w, err_calc.Error())
			return
		}
		task.Date = nextDate
		resultDB := database.DB.Db.Save(&task)
		if resultDB.Error != nil {
			ErrorJson(w, resultDB.Error.Error())
			return
		}
	}

	json_resp, _ := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json_resp)
}
