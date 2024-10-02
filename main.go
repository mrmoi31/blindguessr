package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	game := NewGame()

	var counter int64 = 0

	router := mux.NewRouter()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		counter++
		Connect(w, r, "Player"+strconv.FormatInt(counter, 10), game)
	})

	router.Handle("/", http.FileServer(http.Dir("./")))

	http.ListenAndServe(":8080", router)

	log.Default().Println("SERVER ONLINE")
}
