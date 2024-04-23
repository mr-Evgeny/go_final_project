package model

import (
	"encoding/json"
	"strconv"
)

type ToDo struct {
	ID      uint
	Date    string `gorm:"index,varchar"` //type:datetime
	Title   string `gorm:"varchar"`
	Comment string `gorm:"varchar"`
	Repeat  string `gorm:"type:varchar(128);"`
}

func (ToDo) TableName() string {
	return "scheduler"
}

//fix ID as string for frontend

type jsonToDo struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (p *ToDo) MarshalJSON() ([]byte, error) {
	res := jsonToDo{
		strconv.FormatUint(uint64(p.ID), 10),
		p.Date,
		p.Title,
		p.Comment,
		p.Repeat,
	}

	return json.Marshal(&res)
}
func (p *ToDo) UnmarshalJSON(data []byte) error {
	var res jsonToDo

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	if len(res.ID) > 0 {
		id, err := strconv.ParseUint(res.ID, 10, 32)
		if err != nil {
			return err
		}
		p.ID = uint(id)
	}
	p.Date = res.Date
	p.Title = res.Title
	p.Comment = res.Comment
	p.Repeat = res.Repeat

	return nil
}
