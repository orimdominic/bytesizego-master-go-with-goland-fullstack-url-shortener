package main

import (
	"bytesizego-url-shortener/internal/controllers"
	"bytesizego-url-shortener/internal/db"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	sqlite, err := db.Connect("db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer sqlite.Close()

	err = db.CreateUrlTable(sqlite)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", controllers.Index)
	http.Handle("POST /shorten", controllers.Shorten(sqlite))
	http.Handle("GET /{hash}", controllers.GetURL(sqlite))

	fmt.Println("Listening on 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
