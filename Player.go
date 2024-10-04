package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Player struct {
	name       string
	score      int
	room       *Room
	connection *websocket.Conn
	channel    *Channel
	read       chan string
	write      chan []byte
}

func (p *Player) writePump() {
	for m := range p.write {
		log.Default().Println("PLAYER ", p.name, " WROTE ", string(m))
		p.connection.WriteMessage(websocket.TextMessage, m)
	}
}

func (p *Player) readPump() {
	defer func() {
		p.connection.Close()
		p.room.unregister(p)

	}()
	p.connection.SetReadLimit(1024)

	for {
		_, message, err := p.connection.ReadMessage()
		if err != nil {
			log.Default().Println("Connection closed: ", err)
			break
		}

		log.Default().Println(string(message))

		messageMap := make(map[string]string)
		json.Unmarshal(message, &messageMap)
		p.read <- messageMap["message"]
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Connect(w http.ResponseWriter, r *http.Request, name string, room *Room) *Player {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	player := &Player{
		name:       name,
		score:      0,
		room:       room,
		connection: conn,
		channel:    NewChannel("Private " + name),
		read:       make(chan string),
		write:      make(chan []byte),
	}

	player.channel.players[player] = true

	go player.writePump()
	go player.readPump()

	log.Default().Println("PLAYER CONNECTED : ", name)

	player.room.register(player)

	return player
}
