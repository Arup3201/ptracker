package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/ptracker/internal/domain"
	"github.com/stretchr/testify/assert"
)

var (
	USER_ONE = "ABC128"
	USER_TWO = "BDC920"
)

func mockAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user_id")
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", cookie.Value)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func TestWebsocketNotify(t *testing.T) {
	notifier := NewWsNotifier()
	go notifier.Run()
	handler := NewNotificationHandler([]string{}, notifier)
	server := httptest.NewServer(mockAuth(handler))
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http")
	h := http.Header{}
	h.Add("Cookie", fmt.Sprintf("user_id=%s", USER_ONE))
	client, _, err := websocket.DefaultDialer.Dial(wsUrl, h)
	if err != nil {
		t.Error(err)
		return
	}
	defer client.Close()

	handler.notifier.Notify(context.Background(), USER_ONE, domain.Message{
		Type: "join",
		Data: map[string]string{
			"project_id":   "TABC1",
			"project_name": "Test",
		},
	})

	_, p, err := client.ReadMessage()
	if err != nil {
		t.Error(err)
		return
	}

	var msg domain.Message
	if err := json.Unmarshal(p, &msg); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "join", msg.Type)
	assert.Equal(t, "TABC1", msg.Data["project_id"])
	assert.Equal(t, "Test", msg.Data["project_name"])
}

func TestWebsocketBatchNotify(t *testing.T) {
	notifier := NewWsNotifier()
	go notifier.Run()
	handler := NewNotificationHandler([]string{}, notifier)
	server := httptest.NewServer(mockAuth(handler))
	defer server.Close()

	wsUrl := "ws" + strings.TrimPrefix(server.URL, "http")
	h1 := http.Header{}
	h1.Set("Cookie", fmt.Sprintf("user_id=%s", USER_ONE))
	userOne, _, err := websocket.DefaultDialer.Dial(wsUrl, h1)
	if err != nil {
		t.Error(err)
		return
	}
	defer userOne.Close()

	h2 := http.Header{}
	h2.Set("Cookie", fmt.Sprintf("user_id=%s", USER_TWO))
	userTwo, _, err := websocket.DefaultDialer.Dial(wsUrl, h2)
	if err != nil {
		t.Error(err)
		return
	}
	defer userTwo.Close()

	handler.notifier.BatchNotify(context.Background(),
		[]string{USER_ONE, USER_TWO},
		domain.Message{
			Type: "task_update",
			Data: map[string]string{
				"project_id": "PABC1",
				"task_id":    "TABC1",
				"task_title": "Test",
			},
		})

	_, p, err := userOne.ReadMessage()
	if err != nil {
		t.Error(err)
		return
	}

	var msg domain.Message
	if err := json.Unmarshal(p, &msg); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "task_update", msg.Type)
	assert.Equal(t, "PABC1", msg.Data["project_id"])
	assert.Equal(t, "TABC1", msg.Data["task_id"])
	assert.Equal(t, "Test", msg.Data["task_title"])

	_, p, err = userTwo.ReadMessage()
	if err != nil {
		t.Error(err)
		return
	}
	if err := json.Unmarshal(p, &msg); err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, "task_update", msg.Type)
	assert.Equal(t, "PABC1", msg.Data["project_id"])
	assert.Equal(t, "TABC1", msg.Data["task_id"])
	assert.Equal(t, "Test", msg.Data["task_title"])
}
