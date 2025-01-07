package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients []Client
	channel chan []byte
}

type Client struct {
	server *Server
	conn   *websocket.Conn
}

type Message struct {
	Name string
	Body string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	server := &Server{clients: make([]Client, 0), channel: make(chan []byte)}
	go server.run()

	http.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		client := Client{server: server, conn: conn}
		server.clients = append(server.clients, client)

		fmt.Println("Status: New client connected")

		// Spinning up a new routine for each new client.
		go client.run()
	})

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "app.html")
	})

	http.ListenAndServe(":8080", nil)
}

func (s *Server) run() {
	for {
		msg := <-s.channel

		var message Message
		err := json.Unmarshal(msg, &message)

		if err != nil {
			return
		}

		for _, client := range s.clients {
			if err := client.conn.WriteMessage(1, []byte(fmt.Sprintf("%s: %s", message.Name, message.Body))); err != nil {
				continue
			}
		}
	}
}

func (c *Client) run() {
	defer func() {
		c.conn.Close()
		fmt.Println("Status: Client disconnected")
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		c.server.channel <- msg
	}
}
