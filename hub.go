// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Channel for receiving messages from Redis
	incomingFromRedis chan []byte

	// Redis Client
	redisClient *redis.Client // our connection to the global brain
}

// new Chat-message struct for mapping the incoming JSON data to GO property
// chat-message struct
type ChatMessage struct{
	Type string `json:"type"`
	SenderID string `json:"sender_id"`
	TargetID string `json:"target_id"`
	Content string `json:"content"`

} 

func newHub(rdb *redis.Client) *Hub {
	return &Hub{
		broadcast:         make(chan []byte),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		clients:           make(map[*Client]bool),
		incomingFromRedis: make(chan []byte),
		redisClient:       rdb,
	}
}

func (h *Hub) run() {
	ctx := context.Background()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Create a new empty instance of the ChatMessage struct
			var parsedMsg ChatMessage

			// UNmarshal the incoming json data into the ChatMessage struct 
			err := json.Unmarshal(message, &parsedMsg)
			if err != nil {
				fmt.Println("Error:- ", err)
			}

			// Creating a new channel based on the target_id of the incoming message
			channelName := "user_channel_" + parsedMsg.TargetID

			//publish the message to the Redis based on the selected channel
			err = h.redisClient.Publish(ctx, channelName, message).Err()
			if err != nil {
				fmt.Printf("Error publishing message to Redis: %v", err)
			}
			

		case message := <-h.incomingFromRedis:
			// Subscribe: Redis sent the message! Bradcasting to the local clients

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// listenToRedis runs in the background and constantly listens for new messages on the channel
func (h *Hub) listenToRedis() {
	ctx := context.Background()
	pubsub := h.redisClient.Subscribe(ctx, "chat_room")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		// When a message arrives from Redis, send it to the Hub's internal channel
		h.incomingFromRedis <- []byte(msg.Payload)
	}
}
