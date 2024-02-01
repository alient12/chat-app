package websocket

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/infra/http/request"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Connection struct {
	WS       *websocket.Conn
	Send     chan *model.Message
	IsOnline bool
	UserID   uint64
}

// readPump pumps messages from the websocket connection to the hub.
func (conn *Connection) readPump(c echo.Context, wsc *WebSocketConnection) error {
	defer func() {
		conn.WS.Close()
	}()

	for {
		_, p, err := conn.WS.ReadMessage()
		if err != nil {
			return err
		}

		var req request.MessageCreate

		err = json.Unmarshal(p, &req)
		if err != nil {
			log.Println("cannot unmarshal JSON")
			return echo.ErrBadRequest
		}

		if err := req.Validate(); err != nil {
			log.Print("cannot validate")
			return echo.ErrBadRequest
		}

		msgPtr, err := wsc.mshand.Create(c, req, conn.UserID)
		if err != nil {
			return err
		}

		// send confirmation to the sender
		conn.Send <- msgPtr

		// Send the message to the receiver.
		wsc.lock.RLock()
		receiverConn, ok := wsc.connections[msgPtr.Receiver]
		wsc.lock.RUnlock()
		if ok {
			receiverConn.Send <- msgPtr
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (conn *Connection) writePump(c echo.Context, wsc *WebSocketConnection) {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		conn.WS.Close()
		conn.IsOnline = false
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			conn.WS.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// The hub closed the channel.
				conn.WS.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := conn.WS.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			msg, err := json.Marshal(message)
			if err != nil {
				log.Println("cannot marshal JSON")
				break
			}

			w.Write(msg)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			conn.WS.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WS.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
