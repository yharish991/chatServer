package api

import (
	"errors"

	"github.com/stretchr/testify/mock"

	"chatServer/src/chatserver/data"
)

// ServiceMock struct for api Service struct
type ServiceMock struct {
	mock.Mock
	error error
}


// PostMessage mocks the Service PostMessage method
func (mock *ServiceMock) PostMessage(message data.Message) (data.Message, error) {

	args := mock.Called(message)

	if args.Get(0).(data.Message) != (data.Message{}) {
		return args.Get(0).(data.Message), nil
	}
	return args.Get(0).(data.Message), errors.New("")
}


// GetMessages mocks the Service GetMessages method
func (mock *ServiceMock) GetMessages(userID int, roomID int) (msgs []data.Message) {

	args := mock.Called(userID, roomID)

	if args.Get(0) != nil {
		msgs = args.Get(0).([]data.Message)
	}
	return
}