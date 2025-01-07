package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/listen", func(http.ResponseWriter, *http.Request) {})
	http.ListenAndServe(":8080", nil)
}
