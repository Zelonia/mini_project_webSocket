// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients     = make(map[*websocket.Conn]bool)
	broadcast   = make(chan Message)
	mongoClient *mongo.Client
)

// Message struct untuk menyimpan pesan dan pengirim
type Message struct {
	Username string    `json:"username"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}

func initMongoDB() error {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}
	mongoClient = client
	fmt.Println("Connected to MongoDB!")
	return nil
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, ws)
			break
		}
		msg.Time = time.Now()
		broadcast <- msg

		// Simpan pesan ke MongoDB
		go saveMessageToMongo(msg)
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func saveMessageToMongo(msg Message) {
	collection := mongoClient.Database("chat").Collection("messages")
	_, err := collection.InsertOne(context.TODO(), msg)
	if err != nil {
		fmt.Println(err)
	}
}

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleConnections)
}

func main() {
	err := initMongoDB()
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}

	setupRoutes()

	go handleMessages()

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
