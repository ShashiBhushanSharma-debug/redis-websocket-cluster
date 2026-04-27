# redis-websocket-cluster

> ⚠️ **Work in Progress** — This project is under active development. Features and architecture may change.

A distributed, horizontally scalable WebSocket server built in **Go**, powered by **Redis Pub/Sub** and orchestrated with **Docker Compose**.

This project evolves the classic [gorilla/websocket](https://github.com/gorilla/websocket) chat example into a production-ready distributed system — replacing the monolithic in-memory hub with Redis-backed one-to-one message routing across isolated nodes.

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [What Changed from gorilla/websocket](#what-changed-from-gorillawebsocket)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Roadmap](#roadmap)
- [Acknowledgments](#acknowledgments)

---

## Overview

The standard gorilla/websocket hub stores all client connections in local memory and broadcasts messages to every connected client. This works fine for a single node, but breaks the moment you scale horizontally — a user on Node A is completely invisible to Node B.

This project solves that by introducing **Redis as the message backbone** and replacing broadcast-to-all with **direct, ID-based client routing** across nodes.

**Stack:**
- **Language:** Go (Golang)
- **WebSocket:** [gorilla/websocket](https://github.com/gorilla/websocket)
- **Message Broker:** Redis (Pub/Sub)
- **Infrastructure:** Docker & Docker Compose

---

## Architecture

```
  Client A                        Client B
  (Node 8081)                    (Node 8082)
      │                               │
      ▼                               ▼
 ┌─────────┐    Redis Pub/Sub    ┌─────────┐
 │  Go WS  │◄──────────────────►│  Go WS  │
 │  Node A │                    │  Node B │
 └─────────┘                    └─────────┘
      │                               │
      └───────────┬───────────────────┘
                  │
           ┌─────────────┐
           │    Redis     │
           │  (Broker)   │
           └─────────────┘
```

Each Go node:
1. Maintains its own local map of connected WebSocket clients.
2. Publishes outgoing messages to a Redis channel, tagged with the **target client ID**.
3. Subscribes to Redis and listens for messages routed to its local clients.
4. Delivers the message only to the intended recipient — not to everyone.

---

## What Changed from gorilla/websocket

The official [gorilla/websocket chat example](https://github.com/gorilla/websocket/tree/main/examples/chat) uses a **Hub** — a central struct that holds all connections in memory and fans out every message to all clients.

| | gorilla/websocket (original) | This project |
|---|---|---|
| **Connection registry** | In-memory `map` inside Hub | Redis (shared across nodes) |
| **Message routing** | Broadcast to all clients | Route to a specific client by ID |
| **Multi-node support** | ❌ No | ✅ Yes, via Redis Pub/Sub |
| **Infrastructure** | Single binary | Docker Compose multi-node |

### Key Changes Made

- **`hub.go`** — The hub no longer broadcasts to all local clients. Instead, it publishes messages to a Redis channel with a target client ID embedded in the payload.
- **Redis Subscriber Goroutine** — A background goroutine constantly listens on the Redis channel. When a message arrives, it checks the target ID against locally connected clients and delivers it only to the right one.
- **Client ID Assignment** — Each connecting client is assigned a unique ID, which is used as the routing key.
- **Docker Compose** — Two isolated Go nodes (`8081`, `8082`) and a Redis container are wired together, simulating a real multi-node cluster locally.

---

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

### Run the Cluster

```bash
git clone https://github.com/ShashiBhushanSharma-debug/redis-websocket-cluster.git
cd redis-websocket-cluster
docker compose up --build
```

### Test It

Open two browser windows:

| Node | URL |
|------|-----|
| Node A | http://localhost:8081 |
| Node B | http://localhost:8082 |

Send a message from Node A — it will route through Redis and appear on Node B for the intended recipient only.

---

## Project Structure

```
redis-websocket-cluster/
├── main.go            # Entry point, HTTP server & WebSocket upgrade
├── hub.go             # Modified hub with Redis Pub/Sub integration
├── client.go          # WebSocket client — read/write pumps
├── home.html          # Frontend UI
├── Dockerfile         # Go app container
└── docker-compose.yml # Multi-node cluster setup
```

---

## Roadmap

This project is actively being extended. Planned additions:

- [ ] **Load Balancer** — Distribute incoming WebSocket connections across nodes
- [ ] **Authentication / JWT** — Secure WebSocket handshake with token validation
- [ ] **Horizontal Scaling** — True multi-node deployment with dynamic node discovery
- [ ] **Monitoring / Metrics** — Expose Prometheus metrics for connections, message throughput, and Redis latency

---

## Acknowledgments

The base WebSocket upgrade logic and frontend UI were adapted from the official **[gorilla/websocket chat example](https://github.com/gorilla/websocket/tree/main/examples/chat)** by the [Gorilla Web Toolkit](https://github.com/gorilla) authors.

The backend architecture — including Redis Pub/Sub integration, one-to-one client routing by ID, and Docker multi-node composition — was designed and implemented by **[Shashi Bhushan Sharma](https://github.com/ShashiBhushanSharma-debug)**.

---

*Built with Go · Redis · Docker*