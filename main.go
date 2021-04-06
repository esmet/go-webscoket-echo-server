package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.SetOutput(os.Stdout)
	log.Printf("Echo server listening on port %s.\n", port)
	err := http.ListenAndServe(":"+port, http.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool {
		return true
	},
}

func handler(wr http.ResponseWriter, req *http.Request) {
	log.Printf("%s | %s %s\n", req.RemoteAddr, req.Method, req.URL)
	if websocket.IsWebSocketUpgrade(req) {
		serveWebSocket(wr, req)
	} else if req.URL.Scheme == "ws" {
		wr.Header().Add("Content-Type", "text/html")
		wr.WriteHeader(200)
		io.WriteString(wr, websocketHTML)
	} else {
		serveHTTP(wr, req)
	}
}

func serveWebSocket(wr http.ResponseWriter, req *http.Request) {
	connection, err := upgrader.Upgrade(wr, req, nil)
	if err != nil {
		log.Printf("%s | %s\n", req.RemoteAddr, err)
		return
	}

	defer connection.Close()
	log.Printf("%s | upgraded to websocket\n", req.RemoteAddr)

	var message []byte
	var messageType int

	for {
		// Wait at most 30 seconds for any message...
		connection.SetReadDeadline(time.Now().Add(30 * time.Second))

		messageType, message, err = connection.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			log.Printf("%s | txt | %s\n", req.RemoteAddr, message)
		} else {
			log.Printf("%s | bin | %d byte(s)\n", req.RemoteAddr, len(message))
		}

		err = connection.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}

	if err != nil {
		log.Printf("%s | %s\n", req.RemoteAddr, err)
	}
}

func serveHTTP(wr http.ResponseWriter, req *http.Request) {
	wr.Header().Add("Content-Type", "text/plain")
	wr.WriteHeader(http.StatusNotFound)
}
