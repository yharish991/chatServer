package api

import (
	"errors"

	"chatServer/src/chatserver"
	"chatServer/src/chatserver/data"
)


// ServiceImpl struct for api service
type ServiceImpl struct {
	chatService chatserver.Service
}


// NewServiceImpl returns ServiceImpl
func NewServiceImpl(chatService chatserver.Service, ) *ServiceImpl {
	return &ServiceImpl{
		chatService: chatService,
	}
}


// PostMessage service is for posting a message
func (service *ServiceImpl) PostMessage(message data.Message) (data.Message, error) {

	//validate userID and roomID
	_, userOk := service.chatService.GetUser(message.UserID)
	_, roomOk := service.chatService.GetRoom(message.RoomID)
	if message.UserID != 0 && !userOk {
		return data.Message{}, errors.New("User not found")
	}

	if message.RoomID != 0 && !roomOk {
		return data.Message{}, errors.New("Room not found")
	}

	messageResponse := service.chatService.Publish(data.Input{
		Room: message.RoomID,
		Text:	message.Text,
	}, message.UserID, false)
	return messageResponse, nil
}


// GetMessages service is for retrieving messages
func (service *ServiceImpl) GetMessages(userID int, roomID int) []data.Message {
	messages := service.chatService.GetMessages()

	// if userId and roomId are not provided return all messages
	if userID == 9223372036854775807 && roomID == 9223372036854775807 {
		return messages
	}

	filteredMessages := []data.Message{}
	// filter messages based on the query parameters
	for index, msg := range messages {
		if msg.UserID == userID  && roomID == 9223372036854775807 { // if only userId is provided
			filteredMessages = append(filteredMessages, messages[index])
		} else if msg.RoomID == roomID  && userID == 9223372036854775807 { // if only roomId is provided
			filteredMessages = append(filteredMessages, messages[index])
		} else if msg.RoomID == roomID && msg.UserID == userID { // if both userId and roomId are provided
			filteredMessages = append(filteredMessages, messages[index])
		}
	}
	return filteredMessages
}
