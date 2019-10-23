package api

import "chatServer/src/chatserver/data"

// Service interface for api
type Service interface {
	PostMessage(message data.Message) (data.Message, error)
	GetMessages(userID int, roomID int) []data.Message
}
