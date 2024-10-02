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
	game       *Game
	connection *websocket.Conn
	read       chan string
	write      chan []byte
}

func (p *Player) writePump() {
	for m := range p.write {
		p.connection.WriteMessage(websocket.TextMessage, m)
		log.Default().Println("PLAYER ", p.name, " WROTE ", string(m))
	}
}

func (p *Player) readPump() {
	defer func() {
		p.connection.Close()
		p.game.unregister(p)

	}()
	p.connection.SetReadLimit(1024)

	for {
		_, message, err := p.connection.ReadMessage()
		if err != nil {
			log.Default().Println(err)
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

func Connect(w http.ResponseWriter, r *http.Request, name string, game *Game) *Player {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	player := &Player{
		name:       name,
		score:      0,
		game:       game,
		connection: conn,
		read:       make(chan string),
		write:      make(chan []byte),
	}

	log.Default().Println("PLAYER CONNECTED : ", name)

	go player.writePump()
	go player.readPump()
	
	player.game.register(player)

	return player
}
