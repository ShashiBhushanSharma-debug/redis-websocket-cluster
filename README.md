# Chat Example

This application shows how to use the
[websocket](https://github.com/gorilla/websocket) package to implement a simple
web chat application.

# Distributed Go WebSocket Chat

A horizontally scalable, distributed real-time chat server built with Go, Redis, and Docker. 

This project takes the standard monolithic WebSocket chat architecture and fundamentally re-architects it to solve the **state isolation problem**. By implementing a Redis Pub/Sub message broker, multiple isolated Go nodes can seamlessly communicate, allowing users connected to entirely different containers to chat in real-time.

## 🏗️ Architecture
* **Language:** Go (Golang)
* **WebSocket Handling:** [Gorilla WebSocket](https://github.com/gorilla/websocket)
* **Message Broker:** Redis (Pub/Sub)
* **Infrastructure:** Docker & Docker Compose

## 🚀 The Engineering Challenge: State Isolation
In a standard single-node WebSocket server, the "Hub" stores all connections in local memory. If you spin up two servers behind a load balancer, a user connected to Node A cannot see messages from a user connected to Node B. 

**The Solution:**
Instead of broadcasting directly to local clients, the Go Hub was rewritten to:
1. **Publish:** Intercept incoming messages and immediately push them to a global Redis channel.
2. **Subscribe:** Run a background Goroutine that constantly listens to the Redis channel. When a message is detected, it is pulled down and broadcasted to the local clients.

## ⚙️ How to Run Locally
Ensure you have Docker and Docker Compose installed.

1. Clone the repository.
2. Run the cluster:
   ```bash
   docker compose up --build
3. Open two separate browser windows:

   --> Node A: http://localhost:8081

   --> Node B: http://localhost:8082

4. Type a message in Node A and watch it instantly route through Redis and appear on Node B.


# 🙏 **Acknowledgments**

**The base frontend UI and raw WebSocket connection upgrading logic were adapted from the official Gorilla WebSocket Chat Example. The backend architecture was heavily modified to support distributed Redis Pub/Sub messaging and Docker containerization.**

