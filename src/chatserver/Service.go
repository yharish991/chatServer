package chatserver

import "chatServer/src/chatserver/data"

// Service interface for the chatserver
type Service interface {
	Run()
	CreateUser(username string) data.User
	Publish(input data.Input, userID int, sysMessage bool) data.Message
	Subscribe(userID int, roomID int)
	UnSubscribe(userID int, roomID int)
	SwitchRoom(userID int, roomID int)
	GetActiveRoom(userID int)
	ListRooms(userID int)
	CreateRoom(roomName string, userID int, userName string)
	GetUser(userID int) (data.User, bool)
	GetRoom(roomID int) (data.Room, bool)
	GetMessages() []data.Message
	GetUsers() []data.User
	GetRooms() []data.Room
	RemoveUser(userID int)
}
