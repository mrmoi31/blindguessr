package main

import "math/rand"

type Game struct {
	word     string
	room     *Room
	finished chan bool
}

func StartGame(room *Room) {
	game := &Game{
		word:     Words[rand.Intn(len(Words))],
		room:     room,
		finished: make(chan bool, 1),
	}
	room.play(game)
}
