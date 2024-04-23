package config

import (
	"os"
	"strconv"
)

var (
	SERVER_PORT    string
	SERVER_FILEDIR string
	DB_FILE        string
	PASSWORD       string
	SECRETKEY      string = "top-secret"
)

func Init() {
	port, ok := os.LookupEnv("TODO_PORT")
	if ok && len(port) > 0 {
		if _, err := strconv.ParseInt(port, 10, 16); err == nil {
			SERVER_PORT = port
		}
	}
	if len(SERVER_PORT) == 0 {
		SERVER_PORT = "7540"
	}

	SERVER_FILEDIR = "./web"

	if dbfile, ok := os.LookupEnv("TODO_DBFILE"); ok && len(dbfile) > 0 {
		DB_FILE = dbfile
	} else {
		DB_FILE = "scheduler.db"
	}

	if pass, ok := os.LookupEnv("TODO_PASSWORD"); ok && len(pass) > 0 {
		PASSWORD = pass
	}
}
