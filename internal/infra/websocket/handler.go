package websocket

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/infra/http/handler"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == "http://localhost:3000"
	},
}

type WebSocketConnection struct {
	mshand      *handler.Message
	connections map[uint64]*Connection
	lock        sync.RWMutex
}

func NewWebSocketConnection(mshand *handler.Message) *WebSocketConnection {
	return &WebSocketConnection{
		mshand:      mshand,
		connections: make(map[uint64]*Connection),
		lock:        sync.RWMutex{},
	}
}

func (wsc *WebSocketConnection) WSHandler(c echo.Context) error {
	var uid uint64

	// check auth
	token := c.QueryParam("token")
	if token != "" {
		// check auth by query params
		if ckID, _, err := handler.CheckJWTLocalStorage(token); err != nil {
			return err
		} else {
			uid = ckID
		}
	} else {
		// check auth by cookies
		if ckID, _, err := handler.CheckJWT(c); err != nil {
			return err
		} else {
			uid = ckID
		}
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	conn := &Connection{
		WS:       ws,
		Send:     make(chan *model.Message),
		IsOnline: make(chan bool),
		UserID:   uid,
	}

	// Add the connection to the map.
	wsc.lock.Lock()
	wsc.connections[uid] = conn
	wsc.lock.Unlock()

	// handle connection
	go conn.readPump(c, wsc)
	go conn.writePump(c, wsc)

	return nil
}

func (wsc *WebSocketConnection) Register(g *echo.Group) {
	g.GET("/message", wsc.WSHandler)
}
