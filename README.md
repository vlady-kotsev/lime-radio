# <img width="200px" src="https://github.com/user-attachments/assets/6bd4fad0-6af0-48a8-8cc6-c9ffda368d52" /> Lime Radio 


An internet radio streaming platform with real-time chalga song requests, built with computer technologia.

## Features

- **Real-time MP3 streaming**
- **Song requests** with Solana blockchain payment integration

## Architecture

Lime Radio uses an event-driven microservices architecture with two main services:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│      Client     │    │   Transmitter   │    │       DJ        │
│   (Frontend)    │    │   (Streaming)   │    │  (Song Requests)│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                    │
         │ HTTP/WebSocket        │     Broker         │
         │                       │                    │
         └───────────────────────┼────────────────────┘
                                 │
                    ┌─────────────────┐
                    │     Messaging   │
                    │     Broker      │
                    └─────────────────┘
```

### Transmitter Service (Port 8080)

The core streaming engine responsible for:

- **MP3 Decoding**: Converts MP3 files to PCM audio
- **Smart Buffering**: 200ms buffer-ahead prevents browser audio starvation
- **Precise Timing**: 20ms chunks with monotonic clock for drift-free playback
- **Queue Management**: Real-time song request processing

### DJ Service (Port 8081)

Handles song requests and payments:

- **Solana Integration**: Blockchain payment processing for song requests
- **Event Publishing**: Sends validated requests to transmitter
- **Payment Middleware**: Secures endpoints with transaction verification

### NATS Messaging (Port 4222)

Event-driven communication:

- **JetStream**: Persistent messaging with guaranteed delivery
- **Song Requests**: Real-time request propagation
- **Event Sourcing**: Audit trail for all system events

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Git

### Docker Deployment (Recommended)

1. **Clone the repository:**

   ```bash
   git clone <repository-url>
   cd lime-radio
   ```

2. **Add your MP3 files:**

   ```bash
   # Place MP3 files in the transmitter/songs directory
   cp /path/to/your/songs/*.mp3 transmitter/songs/
   ```

3. **Start the services:**

   ```bash
   docker compose up -d
   ```

## API Reference

### Transmitter Service (8080)

#### Streaming Endpoints

```http
GET /stream
# Returns: WAV audio stream (Content-Type: audio/wav)
# Headers: Connection upgrades to streaming WebSocket
```

#### Song Management

```http
GET /api/songs?page=1&page_size=20&search=artist
# Returns: Paginated song list with metadata

POST /api/songs/update
# Rescans song directory and updates database
```

#### Queue Management

```http
GET /api/queue
# Returns: Currently queued songs

GET /api/queue/count
# Returns: Number of songs in queue
```

### DJ Service (8081)

#### Song Requests

```http
POST /api/request
Content-Type: application/json
{
  "song_id": "uuid",
  "transaction_signature": "solana_signature"
}
# Returns: 200 OK on successful request
```

#### Authentication

```http
POST /api/token
{
  "origin": "http://localhost:3000"
}
# Returns: JWT token for authenticated requests
```
