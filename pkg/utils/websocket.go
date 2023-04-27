package utils

import (
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSLogger struct {
	Conn *websocket.Conn
}

// httpRequest 升级成websocket
func NewWSLogger(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*WSLogger, error) {
	conn, err := upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		return nil, err
	}

	wsConn := &WSLogger{Conn: conn}
	return wsConn, nil
}

func (ws *WSLogger) Write(data []byte) error {
	return ws.Conn.WriteMessage(websocket.TextMessage, data)
}

func (ws *WSLogger) Close() error {
	return ws.Conn.Close()
}
