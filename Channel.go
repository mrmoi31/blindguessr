package main

import (
	"bytes"
	"html/template"
	"log"
)

type Channel struct {
	name    string
	write   chan Message
	players map[*Player]bool
}

type Message struct {
	User    string
	Message string
}

var templ, _ = template.New("message.html").ParseFiles("message.html")

func (c Channel) broadcast() {
	for message := range c.write {
		for player, active := range c.players {
			if active {
				buffer := bytes.Buffer{}
				err := templ.Execute(&buffer, message)
				if err != nil {
					log.Default().Println("Erreur template ", err)
					return
				}
				log.Default().Println("Wrote for " + player.name)
				player.write <- buffer.Bytes()
			}
		}
	}
}

func NewChannel(name string) *Channel {
	channel := &Channel{
		name:    name,
		write:   make(chan Message),
		players: make(map[*Player]bool),
	}

	go channel.broadcast()

	return channel
}
