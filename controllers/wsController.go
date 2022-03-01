package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	helper "github.com/Mutay1/chat-backend/helpers"
	"github.com/Mutay1/chat-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
)

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[string][]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
	UUID   uuid.UUID
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string][]*Client),
}

func updateMessage(message models.Message) (success bool, err error) {
	var friendship models.Friendship
	senderID, _ := primitive.ObjectIDFromHex(message.Sender)
	recipientID, _ := primitive.ObjectIDFromHex(message.RecipientID)

	filter := bson.D{
		{"$or",
			bson.A{
				bson.M{
					"$and": []interface{}{
						bson.M{"requester._id": senderID},
						bson.M{"recipient._id": recipientID},
					},
				},
				bson.M{
					"$and": []interface{}{
						bson.M{"requester._id": recipientID},
						bson.M{"recipient._id": senderID},
					},
				},
			},
		},
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	friendshipCollection.FindOne(ctx, filter).Decode(&friendship)
	if message.UpdateType == "delivered" {
		undelivered := false
		for key, messages := range friendship.Recipient.Messages {
			if !messages.Delivered && messages.RecipientID == message.Sender {
				undelivered = true
				friendship.Recipient.Messages[key].Delivered = true
				friendship.Requester.Messages[key].Delivered = true
			}
		}
		if undelivered {
			data, err := helper.ToDoc(friendship)
			if err != nil {
				return false, err
			}
			_, err = friendshipCollection.UpdateOne(ctx, filter, bson.D{
				{"$set", data},
			})
			if err != nil {
				return false, err
			}
		}
		return true, err
	} else if message.UpdateType == "read" {
		unread := false
		for key, messages := range friendship.Recipient.Messages {
			if !messages.Read && messages.RecipientID == message.Sender {
				unread = true
				friendship.Recipient.Messages[key].Delivered = true
				friendship.Requester.Messages[key].Delivered = true
				friendship.Recipient.Messages[key].Read = true
				friendship.Requester.Messages[key].Read = true
			}
		}
		if unread {
			data, err := helper.ToDoc(friendship)
			if err != nil {
				return false, err
			}
			_, err = friendshipCollection.UpdateOne(ctx, filter, bson.D{
				{"$set", data},
			})
			if err != nil {
				return false, err
			}
		}

	}
	defer cancel()
	return true, err
}

func saveMessage(MessageStruct models.Message) {
	var friendship models.Friendship
	senderID, _ := primitive.ObjectIDFromHex(MessageStruct.Sender)
	recipientID, _ := primitive.ObjectIDFromHex(MessageStruct.RecipientID)

	filter := bson.D{
		{"$or",
			bson.A{
				bson.M{
					"$and": []interface{}{
						bson.M{"requester._id": senderID},
						bson.M{"recipient._id": recipientID},
					},
				},
				bson.M{
					"$and": []interface{}{
						bson.M{"requester._id": recipientID},
						bson.M{"recipient._id": senderID},
					},
				},
			},
		},
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	friendshipCollection.FindOne(ctx, filter).Decode(&friendship)
	friendship.Requester.Messages = append(friendship.Requester.Messages, MessageStruct)
	friendship.Recipient.Messages = append(friendship.Recipient.Messages, MessageStruct)
	data, err := helper.ToDoc(friendship)
	_, err = friendshipCollection.UpdateOne(ctx, filter, bson.D{
		{"$set", data},
	})
	if err != nil {
		fmt.Println(err)
	}
	defer cancel()
}

func remove(s []*Client, i int) []*Client {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

//Start is before the project runs, the program starts start > go Manager.Start ()
func (manager *ClientManager) Start() {
	for {
		log.Println("< --- pipeline communication -- >")
		select {
		case conn := <-Manager.Register:
			log.Printf(("new user joined in% v"), conn.ID)
			Manager.Clients[conn.ID] = append(Manager.Clients[conn.ID], conn)
			jsonMessage, _ := json.Marshal(&models.Message{Content: "Successful connection to socket service"})
			conn.Send <- jsonMessage
		case conn := <-Manager.Unregister:
			log.Printf(("user left% v"), conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				// jsonMessage, _ := json.Marshal(&models.Message{Content: "A socket has disconnected"})
				// conn.Send <- jsonMessage
				// if len(Manager.Clients[conn.ID]) <= 1 {
				// 	close(conn.Send)
				// }
				for index, c := range manager.Clients[conn.ID] {
					if c.UUID == conn.UUID {
						manager.Clients[conn.ID] = remove(manager.Clients[conn.ID], index)
					}
				}
			}
		case message := <-Manager.Broadcast:
			MessageStruct := models.Message{}
			json.Unmarshal(message, &MessageStruct)

			fmt.Println(MessageStruct)
			if MessageStruct.MessageType == "info" {
				success, _ := updateMessage(MessageStruct)
				if success {
					for id, conns := range Manager.Clients {
						if id == MessageStruct.Sender || id == MessageStruct.RecipientID {
							for _, conn := range conns {
								conn.Send <- message
							}
						}
					}
				}
			} else {
				for id, conns := range Manager.Clients {
					if id == MessageStruct.Sender || id == MessageStruct.RecipientID {
						for _, conn := range conns {
							conn.Send <- message
						}
					}
				}
				saveMessage(MessageStruct)
			}
		}
	}
}
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		log.Printf("message read to client: s", string(message))
		Manager.Broadcast <- message
	}
}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("message sent to client: s", string(message))
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}

	}
}

// func (c *Client) PingHandler() {
// 	defer func() {
// 		c.Socket.Close()
// 	}()
// 	MessageStruct := models.Message{
// 		MessageType: "info",
// 		Content:     "Ping",
// 	}
// 	message, err := json.Marshal(MessageStruct)
// 	if err != nil {
// 		Manager.Unregister <- c
// 		c.Socket.Close()
// 	}
// 	for _ = range time.Tick(5 * time.Second) {
// 		log.Printf("message sent to client: s", string(message))
// 		c.Socket.WriteMessage(websocket.TextMessage, message)
// 	}
// }

//WsHandler socket connection middleware function: upgrade protocol, user authentication, user-defined information, etc
func WsHandler(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//User information authentication can be added
	id, _ := uuid.New()
	client := &Client{
		ID:     c.Query("uid"),
		Socket: conn,
		Send:   make(chan []byte),
		UUID:   id,
	}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func Pong() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	}
}
