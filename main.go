package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	LoadWords()

	room := NewRoom()

	var counter int64 = 0

	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		counter++
		Connect(w, r, room)
	})

	//router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	router.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		open, err := os.Open("static/css/styles.css")
		if err != nil {
			return
		}
		w.Header().Set("Content-Type", "text/css")
		open.WriteTo(w)
	})

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		log.Default().Println("Looking for favicon")
		open, err := os.Open("static/favicon.ico")
		if err != nil {
			return
		}
		open.WriteTo(w)
	})

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		open, err := os.Open("html/index.html")
		if err != nil {
			return
		}
		open.WriteTo(w)
	})

	log.Default().Println("SERVER ONLINE")

	http.ListenAndServe(":8080", router)

}
