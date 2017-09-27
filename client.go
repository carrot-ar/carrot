package buddy

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (

	//time allowed to write a message to the websocket
	writeWait = 10 * time.Second

	//time allowed to read the next pong message from the websocket
	pongWait = 10 * time.Second

	//send pings to the websocket with this period, must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	//maximum message size allowed from the websocket
	maxMessageSize = 65536

	// Toggle to require a client secret token on WS upgrade request
	clientSecretRequired = false
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	server *Server
	open   bool

	conn *websocket.Conn

	//buffered channel of outbound messages
	send chan []byte

	//buffered channel of outbound tokens
	sendToken chan SessionToken
}

//readPump pumps messages from the websocket to the server
func (c *Client) readPump() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.server.broadcast <- message
	}
}

//writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//the server closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			//add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case token, ok := <-c.sendToken:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				//the server closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(token))

			//add queued messages to the current websocket message
			n := len(c.sendToken)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write([]byte(<-c.sendToken))
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

//authenticate that the client should be allowed to connect

func validClientSecret(clientSecret string) bool {
	log.Printf("clientSecret: %v", clientSecret)

	if clientSecret == serverSecret {
		return true
	}

	log.Println("The client and server secrets do not match :(")
	return false
}

func serveWs(server *Server, w http.ResponseWriter, r *http.Request) {
	sessionToken, clientSecret, _ := r.BasicAuth()
	log.Printf("Session Token: %v | Client Secret: %v", sessionToken, clientSecret)

	if clientSecretRequired && !validClientSecret(clientSecret) {
		http.Error(w, "Not Authorized!", 403)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{server: server, conn: conn, send: make(chan []byte, 256), sendToken: make(chan SessionToken, 256)}
	client.sendToken <- SessionToken(sessionToken)
	client.server.register <- client

	go client.writePump()
	go client.readPump()
}
