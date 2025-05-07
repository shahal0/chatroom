# Go Gin Chat Application

A simple concurrent chat server built with Go, Gin, and Server-Sent Events (SSE). Users can join, send messages, leave, and receive real-time chat updates. System messages notify when users join or leave.

## Features
- Join and leave chat rooms
- Send and receive messages in real time
- See your own messages as well as others'
- System notifications when users join or leave
- Simple RESTful API endpoints
- Graceful server shutdown

## Requirements
- Go 1.18 or newer
- Gin web framework

## Installation
1. Clone the repository:
   ```sh
   git clone https://github.com/shahal0
   cd chat
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```

## Running the Server
```sh
go run ./cmd/main.go
```
The server will start on port 3000 by default.

## API Endpoints

### Join Chat
- **Endpoint:** `GET /join?id=<client_id>`
- **Description:** Join the chat room with a unique client ID.
- **Response:**
  ```json
  { "success": true, "message": "Client <id> joined the chat room" }
  ```

### Send Message
- **Endpoint:** `GET /send?id=<client_id>&message=<text>`
- **Description:** Send a message as the specified client.
- **Response:**
  ```json
  { "success": true, "message": "Message sent from client <id>" }
  ```

### Leave Chat
- **Endpoint:** `GET /leave?id=<client_id>`
- **Description:** Leave the chat room.
- **Response:**
  ```json
  { "success": true, "message": "Client <id> left the chat room" }
  ```

### Receive Messages (SSE)
- **Endpoint:** `GET /messages?id=<client_id>`
- **Description:** Receive real-time messages as Server-Sent Events. Each message is streamed as:
  ```
  data: <sender_id>:<message>\n\n
  Example: data: 1:hello\n\n
  System messages (join/leave) are sent as:
  data: system:<client_id> joined\n\n
  Timeout after 30 seconds of inactivity:
  data: timeout\n\n
  ```
- **Note:** The client must keep the connection open to receive messages.

## Example Usage

1. **Join:**
   ```sh
   curl "http://localhost:3000/join?id=1"
   ```
2. **Send Message:**
   ```sh
   curl "http://localhost:3000/send?id=1&message=hello"
   ```
3. **Receive Messages:**
   ```sh
   curl "http://localhost:3000/messages?id=1"
   ```
4. **Leave:**
   ```sh
   curl "http://localhost:3000/leave?id=1"
   ```

## Notes
- Each client must use a unique `id`.
- If a client tries to join with an existing ID, an error is returned.
- Messages are only received while connected to `/messages`.
- System messages are sent when users join or leave.
- The server supports graceful shutdown on SIGINT/SIGTERM.

## License
MIT