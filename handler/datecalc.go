package handler

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type DatesArray [3][]int   //date, weekday, month
const DayLimit int = 10000 //inf loop protection, upto 27 years limit

func calcDayShift(date time.Time, dates *DatesArray) int {
	cnt := 0
	for cnt < DayLimit {
		d := date.Day()
		w := int(date.Weekday())
		//us weektime to rus
		if w == 0 {
			w = 7
		}
		m := int(date.Month())
		last_day := date.AddDate(0, 1, -d).Day()
		var mask byte = 0
		for i := 0; i < 3; i++ {
			for _, v := range dates[i] {
				if i == 0 && v < 0 {
					v = last_day + v + 1
				}
				if (i == 0 && v == d) || (i == 1 && v == w) || (i == 2 && v == m) {
					mask |= (1 << i)
					break
				}
			}
			//if array empty set bits up
			if len(dates[i]) == 0 {
				mask |= (1 << i)
			}
		}
		if mask == 7 { //7 eq 3 bits up
			return cnt
		}
		date = date.AddDate(0, 0, 1)
		cnt++
	}
	return cnt
}

func getNextDate(now time.Time, date_str string, repeat string) (result string, err error) {
	if len(repeat) == 0 {
		err = errors.New("field 'repeat' is empty")
		return
	}
	date, errParse := time.Parse("20060102", date_str)
	if errParse != nil {
		err = errors.New("field 'date' has wrong format")
		return
	}
	repeat_arr := strings.Split(repeat, " ")

	//preparse extra repeat info && checking format
	var dates DatesArray //days, weeks, months
	if len(repeat_arr) > 1 && len(repeat_arr) < 4 && repeat_arr[0] == "m" {
		//list of days in month
		for _, day := range strings.Split(repeat_arr[1], ",") {
			if d, errconv := strconv.ParseInt(day, 10, 16); errconv == nil && d >= -2 && d <= 31 {
				dates[0] = append(dates[0], int(d))
			} else {
				err = errors.New("repeat month format is wrong")
				return
			}
		}
		//list of months
		if len(repeat_arr) == 3 {
			for _, month := range strings.Split(repeat_arr[2], ",") {
				if m, errconv := strconv.ParseInt(month, 10, 16); errconv == nil && m > 0 && m <= 12 {
					dates[2] = append(dates[2], int(m))
				} else {
					err = errors.New("repeat month format is wrong")
					return
				}
			}
		}
	} else if len(repeat_arr) == 2 && repeat_arr[0] == "w" {
		//list of weekdays
		for _, wday := range strings.Split(repeat_arr[1], ",") {
			if w, errconv := strconv.ParseInt(wday, 10, 16); errconv == nil && w > 0 && w <= 7 {
				dates[1] = append(dates[1], int(w))
			} else {
				err = errors.New("repeat week format is wrong")
				return
			}
		}
	} else if len(repeat_arr) == 2 && repeat_arr[0] == "d" {
		//day shift
		if d, errconv := strconv.ParseInt(repeat_arr[1], 10, 16); errconv == nil && d > 0 && d <= 400 {
			dates[0] = append(dates[0], int(d))
		} else {
			err = errors.New("repeat day format is wrong")
			return
		}
	} else if len(repeat_arr) != 1 || repeat_arr[0] != "y" {
		err = errors.New("repeat format is wrong")
		return
	}
	//result should greater than current date + fix(can process tasks in future)
	cnt := 0
	for date.Before(now) || cnt == 0 {
		if cnt > DayLimit {
			break
		}
		cnt++
		if repeat_arr[0] == "d" {
			date = date.AddDate(0, 0, dates[0][0])
		} else if repeat_arr[0] == "y" {
			date = date.AddDate(1, 0, 0)
		} else if repeat_arr[0] == "w" || repeat_arr[0] == "m" {
			if date.Before(now) || date.Equal(now) {
				date = now.AddDate(0, 0, 1)
			}
			date = date.AddDate(0, 0, calcDayShift(date, &dates))
		}
	}
	result = date.Format("20060102")
	return
}

func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}
