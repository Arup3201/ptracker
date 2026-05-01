package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ptracker/notifications"
)

type ListedMessages struct {
	Messages []notifications.Notification `json:"messages"`
}

type MessageApi struct {
	notificationService *notifications.NotificationService
}

func NewMessageApi(
	notificationService *notifications.NotificationService,
) *MessageApi {
	return &MessageApi{
		notificationService: notificationService,
	}
}

func (api *MessageApi) List(w http.ResponseWriter, r *http.Request) error {

	userID, err := GetUserID(r)
	if err != nil {
		return fmt.Errorf("get user Id: %w", err)
	}

	notifications, err := api.notificationService.List(
		r.Context(),
		userID,
	)

	json.NewEncoder(w).Encode(HTTPSuccessResponse[ListedMessages]{
		Data: &ListedMessages{
			Messages: notifications,
		},
	})

	return nil
}
