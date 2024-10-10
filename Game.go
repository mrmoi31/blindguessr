package main

import (
	"math"
	"time"
)

type Game struct {
	word     string
	room     *Room
	start    time.Time
	duration int
	finished chan bool
}

func StartGame(room *Room, duration int) {
	game := &Game{
		word:     RandomWord(),
		room:     room,
		start:    time.Now(),
		duration: duration,
		finished: make(chan bool, 1),
	}
	go game.checkOver()
	room.play(game)
}

func (game *Game) RemainingTime() int {
	return game.duration - int(math.Trunc(time.Now().Sub(game.start).Seconds()))
}

func (game *Game) checkOver() {
	t := -1

	for {
		if t != game.RemainingTime() {

			t = game.RemainingTime()

			if t <= 0 {
				game.finished <- false
				return
			}
		}
	}
}
