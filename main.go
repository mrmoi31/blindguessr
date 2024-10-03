package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {

	room := NewRoom()

	var counter int64 = 0

	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		counter++
		Connect(w, r, "Player"+strconv.FormatInt(counter, 10), room)
	})

	router.Handle("/", http.FileServer(http.Dir("./")))

	log.Default().Println("SERVER ONLINE")

	http.ListenAndServe(":8080", router)

}
