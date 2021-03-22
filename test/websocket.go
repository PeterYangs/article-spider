package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func main() {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		if err != nil {

			fmt.Println(err)

			return
		}

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()

			fmt.Println(msgType)

			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, "websockets.html")
	//})

	http.ListenAndServe(":8088", nil)
}

func checkOrigin(r *http.Request) bool {

	return true
}
