package chatserver

import (
	"github.com/stretchr/testify/mock"

	"chatServer/src/chatserver/data"
)

// ServiceMock struct for chat server Service struct
type ServiceMock struct {
	mock.Mock
	messages []data.Message
}


// Run mocks chatserver Service Run method
func (mock *ServiceMock) Run() {
}


// CreateUser mocks chatserver Service CreateUser method
func (mock *ServiceMock) CreateUser(username string) data.User {
	return data.User{}
}


// Publish mocks chatserver Service Publish method
func (mock *ServiceMock) Publish(input data.Input, userID int, sysMessage bool) data.Message {
	return dummyMessages[1]
}


// Subscribe mocks chatserver Service Subscribe method
func (mock *ServiceMock) Subscribe(userID int, roomID int) {
}


// GetActiveRoom mocks chatserver Service GetActiveRoom method
func (mock *ServiceMock) GetActiveRoom(userID int) {
}


// UnSubscribe mocks chatserver Service UnSubscribe method
func (mock *ServiceMock) UnSubscribe(userID int, roomID int) {
}


// SwitchRoom mocks chatserver Service SwitchRoom method
func (mock *ServiceMock) SwitchRoom(userID int, roomID int) {
}


// ListRooms mocks chatserver Service ListRooms method
func (mock *ServiceMock) ListRooms(userID int) {
}


// CreateRoom mocks chatserver Service CreateRoom method
func (mock *ServiceMock) CreateRoom(roomName string, userID int, userName string) {
}


// SendInfo mocks chatserver Service SendInfo method
func (mock *ServiceMock) SendInfo(info string, userID int) {
}


// GetRoom mocks chatserver Service GetRoom method
func (mock *ServiceMock) GetRoom(roomID int) (data.Room, bool) {
	return data.Room{}, false
}


// GetUser mocks chatserver Service GetUser method
func (mock *ServiceMock) GetUser(userID int) (data.User, bool) {
	if userID > len(dummyUsers) {
		return data.User{}, false
	}
	return dummyUsers[userID], true
}

// GetMessages mocks chatserver Service GetMessages method
func (mock *ServiceMock) GetMessages() []data.Message {
	return dummyMessages
}

// GetUsers mocks chatserver Service GetUsers method
func (mock *ServiceMock) GetUsers() []data.User {
	return dummyUsers
}

// GetRooms mocks chatserver Service GetRooms method
func (mock *ServiceMock) GetRooms() []data.Room {
	return dummyRooms
}

// RemoveUser mocks chatserver Service RemoveUser method
func (mock *ServiceMock) RemoveUser(userID int) {
}

var dummyMessages = []data.Message {
	{
		ID: 0,
		UserID: 1,
		RoomID: 0,
		UserName: "harish",
		RoomName: "Default",
		Text: "hello",
		TimeStamp: "20190608172307",
	},
	{
		ID: 1,
		UserID: 2,
		RoomID: 0,
		UserName: "abhilash",
		RoomName: "Default",
		Text: "hi this is abhilash",
		TimeStamp: "20190608172338",
	},
	{
		ID: 2,
		UserID: 3,
		RoomID: 0,
		UserName: "anusha",
		RoomName: "Default",
		Text: "hi this is anusha",
		TimeStamp: "20190608172359",
	},
}

var dummyUsers = []data.User {
	{
		ID: 0,
		Name: "System",
		ActiveRoom: 0,
	},
	{
		ID: 1,
		Name: "Rob",
		ActiveRoom: 0,
	},
	{
		ID: 2,
		Name: "Bob",
		ActiveRoom: 0,
	},
	{
		ID: 3,
		Name: "John",
		ActiveRoom: 0,
	},
}

var dummyRooms = []data.Room {
	{
		ID: 0,
		Name: "Default",
	},
}