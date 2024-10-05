package main

import (
	"bufio"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
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

func Register(conn *websocket.Conn, name string, room *Room) *Player {

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

	landing, err := os.Open("html/game.html")

	if err != nil {
		log.Fatal("Erreur ouverture landing : ", err)
		return nil
	}

	reader := bufio.NewReader(landing)

	bytes := make([]byte, 5000)
	_, err = reader.Read(bytes)
	if err != nil {
		log.Fatal("Erreur lecture landing : ", err)
		return nil
	}

	err = conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Fatal("Erreur write : ", err)
		return nil
	}
	log.Default().Println("Wrote : ", string(bytes))

	player.room.register(player)

	return player
}

func Connect(w http.ResponseWriter, r *http.Request, room *Room) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Erreur upgrade : ", err)
		return
	}

	landing, err := os.Open("html/landing.html")

	if err != nil {
		log.Fatal("Erreur ouverture landing : ", err)
		return
	}

	reader := bufio.NewReader(landing)

	bytes := make([]byte, 1000)
	_, err = reader.Read(bytes)
	if err != nil {
		log.Fatal("Erreur lecture landing : ", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		log.Fatal("Erreur write : ", err)
		return
	}

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Default().Println("Connection closed: ", err)
				break
			}

			log.Default().Println("Username message : ", string(message))

			messageMap := make(map[string]string)
			json.Unmarshal(message, &messageMap)

			username, ok := messageMap["username"]
			if !ok {
				log.Default().Println("No username")
				continue
			}

			Register(conn, username, room)
			break
		}
	}()
}
