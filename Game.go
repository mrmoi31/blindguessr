package main

import "log"

type Game struct {
	players  map[*Player]bool
	channels map[string]*Channel
}

func NewGame() *Game {
	game := &Game{
		players:  make(map[*Player]bool),
		channels: make(map[string]*Channel),
	}

	game.channels["global"] = NewChannel("Global")

	return game
}

func (game *Game) register(player *Player) {
	game.players[player] = true
	go game.readPlayer(player)
	game.channels["global"].players[player] = true
}

func (game *Game) unregister(player *Player) {
	delete(game.players, player)
}

func (game *Game) play() {
	for player := range game.players {
		go game.readPlayer(player)
	}
}

func (game *Game) readPlayer(player *Player) {
	for message := range player.read {
		log.Default().Println("READ ", message, " ON PLAYER ", player.name)

		game.channels["global"].write <- Message{User: player.name, Message: message}
	}
}
