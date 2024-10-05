package main

import "log"

type Room struct {
	players map[*Player]bool
	global  *Channel
	winners *Channel
	game    *Game
}

func NewRoom() *Room {
	room := &Room{
		players: make(map[*Player]bool),
		global:  NewChannel("Global"),
		winners: NewChannel("Winners"),
		game:    nil,
	}

	return room
}

func (room *Room) register(player *Player) {
	room.players[player] = false
	go room.readPlayer(player)
	room.global.players[player] = true

	if len(room.players) == 1 {
		StartGame(room)
	}
	room.global.write <- Message{User: "Game", Message: player.name + " joined us"}
}

func (room *Room) unregister(player *Player) {
	delete(room.global.players, player)
	delete(room.winners.players, player)
	delete(room.players, player)
}

func (room *Room) play(game *Game) {
	if game == nil {
		return
	}
	room.game = game
	for player, _ := range room.players {
		room.players[player] = false
		delete(room.winners.players, player)
	}

	room.global.write <- Message{User: "Game", Message: "New game started!"}
	log.Default().Println("Word is " + game.word)

	over := <-room.game.finished

	if over {
		room.global.write <- Message{User: "Game", Message: "Good job everyone!"}
		room.global.write <- Message{User: "Game", Message: "The word was : " + game.word}
	} else {
		room.global.write <- Message{User: "Game", Message: "Game over!"}
		room.global.write <- Message{User: "Game", Message: "The word was : " + game.word}
	}

	StartGame(room)
}

func (room *Room) readPlayer(player *Player) {
	for message := range player.read {
		log.Default().Println("READ ", message, " ON PLAYER ", player.name)

		if room.players[player] {
			room.winners.write <- Message{User: player.name, Message: message}
		} else if message == room.game.word {
			room.players[player] = true
			room.winners.players[player] = true
			room.global.write <- Message{User: "Game", Message: player.name + " guessed the word!"}
			player.channel.write <- Message{User: "Game", Message: "Good guess, the word was indeed : " + message}

			if len(room.winners.players) == len(room.players) {
				room.game.finished <- true
				close(room.game.finished)
			}
		} else {
			room.global.write <- Message{User: player.name, Message: message}
		}
	}
}
