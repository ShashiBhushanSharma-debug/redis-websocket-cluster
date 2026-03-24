// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./index/home.html")
}

func main() {
	flag.Parse()

	// 1. Grab the Redis URL from Docker Compose environment variables
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379" // Fallback if running outside Docker
	}

	// 2. Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// 3. Initialize the Hub with Redis
	hub := newHub(rdb)

	// 4. Start the background processes
	go hub.listenToRedis() // Listen for global messages
	go hub.run()           // Manage local state

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
