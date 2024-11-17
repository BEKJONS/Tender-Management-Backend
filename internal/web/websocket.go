package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"tender_management/config"
	"tender_management/internal/email"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Production uchun sozlang!
		},
	}
	clients = make(map[*Client]bool)
	ctx     = context.Background()
	rdb     *redis.Client
)

type Client struct {
	conn *websocket.Conn
}

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

func handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocketga ulanish xatosi: %v", err)
		return
	}
	client := &Client{conn: conn}
	clients[client] = true
	fmt.Println("Yangi mijoz ulanishi amalga oshdi")
	go handleMessages(client)
}

func handleMessages(client *Client) {
	defer func() {
		client.conn.Close()
		delete(clients, client)
		fmt.Println("Mijoz uzildi")
	}()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Xabarni o'qishda xato: %v", err)
			break
		}
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Xabarni JSONga aylantirish xatosi: %v", err)
			continue
		}
		if msg.Type == "notification" {
			// Redisga xabar saqlash
			err := rdb.LPush(ctx, "notifications", msg.Content).Err()
			if err != nil {
				log.Printf("Redisga yozishda xato: %v", err)
				continue
			}
			broadcastMessage(message)
		}
	}
}

func broadcastMessage(message []byte) {
	for client := range clients {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Xabarni tarqatishda xato: %v", err)
			client.conn.Close()
			delete(clients, client)
		}
	}
}

func SendNotification(c *gin.Context, message string, Config *config.Config, rdb *redis.Client, Email string) {
	if message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message parameter required"})
		return
	}
	// Redisga xabarni saqlash
	rdb.LPush(ctx, "notifications", message)
	broadcastMessage([]byte(message))
	email.SendEmail(ctx, Config, rdb, Email, message)
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent!"})
}

func WebSocketRunner() {
	initRedis()
	r := gin.Default()
	r.GET("/ws", handleWebSocket)
	fmt.Println("WebSocket server ishlamoqda :8080 portida")
	r.Run(":8080")
}
