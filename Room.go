package main

import "log"

type Room struct {
	players map[*Player]bool
	global  *Channel
}

func NewRoom() *Room {
	room := &Room{
		players: make(map[*Player]bool),
		global:  NewChannel("Global"),
	}

	return room
}

func (room *Room) register(player *Player) {
	room.players[player] = true
	go room.readPlayer(player)
	room.global.players[player] = true
}

func (room *Room) unregister(player *Player) {
	delete(room.players, player)
}

func (room *Room) play() {
	for player := range room.players {
		go room.readPlayer(player)
	}
}

func (room *Room) readPlayer(player *Player) {
	for message := range player.read {
		log.Default().Println("READ ", message, " ON PLAYER ", player.name)
		room.global.write <- Message{User: player.name, Message: message}
	}
}
