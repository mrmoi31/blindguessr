package main

type Game struct {
	word     string
	room     *Room
	finished chan bool
}

func StartGame(room *Room) {
	game := &Game{
		word:     RandomWord(),
		room:     room,
		finished: make(chan bool, 1),
	}
	room.play(game)
}
