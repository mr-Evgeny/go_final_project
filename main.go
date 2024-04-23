package main

import (
	"github.com/mr-Evgeny/go_final_project/config"
	"github.com/mr-Evgeny/go_final_project/database"
	"github.com/mr-Evgeny/go_final_project/handler"
	"log"
	"net/http"
)

func main() {
	config.Init()
	database.Connect()

	fileServer := http.FileServer(http.Dir(config.SERVER_FILEDIR))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/signin", handler.Sign)
	http.HandleFunc("/api/nextdate", handler.NextDate)
	http.HandleFunc("/api/", handler.Auth(handler.Api))

	log.Printf("Starting server at port %s\n", config.SERVER_PORT)
	if err := http.ListenAndServe(":"+config.SERVER_PORT, nil); err != nil {
		log.Fatal(err)
	}
}
