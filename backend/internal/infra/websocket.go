package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/utils"
)

const (
	WRITE_WAIT = 10 * time.Second
	PING_WAIT  = 50 * time.Second
)

type notificationHandler struct {
	upgrader websocket.Upgrader
	notifier *wsNotifier
}

func NewNotificationHandler(allowOrigins []string, n *wsNotifier) http.Handler {

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if len(allowOrigins) == 0 {
				return true
			}

			origin := r.Header.Get("Origin")
			return slices.Contains(allowOrigins, origin)
		},
	}

	return &notificationHandler{
		upgrader: upgrader,
		notifier: n,
	}
}

func (h *notificationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	userId, err := utils.GetUserId(r)
	if err != nil {
		conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
		conn.WriteMessage(websocket.CloseMessage, []byte{})
		conn.Close()
		return
	}

	id := uuid.NewString()
	client := &wsClient{
		conn:     conn,
		notifier: h.notifier,

		id:   id,
		user: userId,

		send: make(chan domain.Message),
	}
	h.notifier.register <- client

	go client.writePump()
}

type wsNotifier struct {
	mutex sync.RWMutex

	register   chan *wsClient
	unregister chan *wsClient

	clients map[string][]*wsClient
}

func NewWsNotifier() *wsNotifier {

	return &wsNotifier{
		mutex: sync.RWMutex{},

		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		clients:    make(map[string][]*wsClient),
	}
}

func (n *wsNotifier) Run() {
	for {
		select {
		case client := <-n.register:
			n.mutex.Lock()

			if _, ok := n.clients[client.user]; !ok {
				n.clients[client.user] = []*wsClient{}
			}

			n.clients[client.user] = append(n.clients[client.user], client)

			n.mutex.Unlock()
		case client := <-n.unregister:
			n.mutex.Lock()

			if _, ok := n.clients[client.user]; ok {
				if len(n.clients[client.user]) == 1 {
					close(client.send)
					delete(n.clients, client.user)
				} else {
					i := slices.IndexFunc(n.clients[client.user], func(c *wsClient) bool {
						return c.id == client.id
					})
					if i == -1 {
						fmt.Println("[ERROR] client not found during unregister")
						continue
					}

					n.clients[client.user] = n.clients[client.user][i : i+1]
				}
			}

			n.mutex.Unlock()
		}
	}
}

func (n *wsNotifier) Notify(ctx context.Context,
	user string, message domain.Message) error {
	clients, ok := n.clients[user]
	if !ok {
		return fmt.Errorf("user is offline")
	}

	for _, client := range clients {
		client.send <- message
	}

	return nil
}

func (n *wsNotifier) BatchNotify(ctx context.Context,
	users []string, message domain.Message) error {

	for _, user := range users {
		err := n.Notify(ctx, user, message)
		if err != nil {
			fmt.Println("[ERROR] batch notify: %w", err)
		}
	}

	return nil
}

type wsClient struct {
	conn     *websocket.Conn
	notifier *wsNotifier

	id   string
	user string

	send chan domain.Message
}

// write to client connection
func (c *wsClient) writePump() {
	ticker := time.NewTicker(PING_WAIT)
	defer func() {
		ticker.Stop()
		c.notifier.unregister <- c
		c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
		c.conn.WriteMessage(websocket.CloseMessage, nil)
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				break
			}

			bytes, err := json.Marshal(message)
			if err != nil {
				continue
			}

			_, err = w.Write(bytes)
			if err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}
	}
}
