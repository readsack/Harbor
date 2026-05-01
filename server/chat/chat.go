package chat

import (
	"fmt"
	"harbor/main/db"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Send chan db.Message
	User *db.User
	Conn *websocket.Conn
	db.Chat
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan db.Message
}

var Upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var centralHub *Hub

func GetHub() *Hub {
	return centralHub
}

func InitHub() {
	centralHub = &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan db.Message),
	}

	go centralHub.HandleClients()
}

func (client *Client) Listen() {
	//fmt.Println("Listener Started")
	hub := GetHub()
	defer func() {
		hub.unregister <- client
		client.Conn.Close()
	}()
	for {
		var m db.Message
		err := client.Conn.ReadJSON(&m)
		//log.Println(m)

		if err != nil {
			log.Println(err)
			return
		}
		m.Name = client.User.Username
		m.ID = client.User.ID
		m.ChatID = client.Chat.ID
		go db.AddMsg(m.Content, m.ChatID, client.User.ID)
		hub.broadcast <- m
	}
}

func (client *Client) SendMessage() {
	//fmt.Println("Send Message Listener Started")
	defer client.Conn.Close()
	for {

		msg := <-client.Send
		fmt.Println(msg)

		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (hub *Hub) HandleClients() {
	//fmt.Println("Central Hub Started")
	//defer fmt.Println("Central Hub Ended")

	for {
		select {
		case msg := <-hub.broadcast:

			for client := range hub.clients {
				if msg.ChatID == client.Chat.ID {
					client.Send <- msg
				}
			}
		case client := <-hub.register:
			hub.clients[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.clients[client]; ok {
				close(client.Send)
				delete(hub.clients, client)
			}
		}
	}

}

func (client *Client) HandleClientConnection() {
	fmt.Println(client.User.Username)
	go client.Listen()
	go client.SendMessage()
	hub := GetHub()
	hub.register <- client
}
