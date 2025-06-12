package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"sensor_dashboard_influx/sensors"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket connection to the hub.
type Client struct {
	// The WebSocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// Reference to the hub.
	Hub *Hub
}

// commandMessage represents messages received from clients
type commandMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// setMotorPayload holds the payload for motor commands
type setMotorPayload struct {
	On bool `json:"on"`
}

var upgrader = websocket.Upgrader{
	// Allow connections from any origin; CORS is handled at HTTP layer.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ServeWS handles WebSocket requests from clients.
func ServeWS(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP to WebSocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("websocket upgrade error: %v", err)
			return
		}

		client := &Client{
			Conn: conn,
			Send: make(chan []byte, 256),
			Hub:  hub,
		}

		hub.Register <- client

		// Start write pump
		go client.writePump()
		// Start read pump
		client.readPump()
	}
}

// readPump reads messages from the WebSocket, handling control commands.
func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// Parse command
		var cmd commandMessage
		if err := json.Unmarshal(msg, &cmd); err != nil {
			continue
		}
		switch cmd.Type {
		case "setMotor":
			var p setMotorPayload
			if err := json.Unmarshal(cmd.Payload, &p); err == nil {
				sensors.SetManualMotor(p.On)
				update := map[string]interface{}{"type": "motorStatus", "payload": map[string]bool{"on": p.On}}
				if data, err := json.Marshal(update); err == nil {
					c.Hub.Broadcast <- data
				}
			}
		case "autoMotor":
			sensors.ClearManualMotor()
		default:
			// ignore unknown commands
		}
	}
}

// writePump writes messages from the hub to the WebSocket.
func (c *Client) writePump() {
	defer c.Conn.Close()
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
