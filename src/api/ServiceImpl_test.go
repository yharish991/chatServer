package api

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"chatServer/src/chatserver"
	"chatServer/src/chatserver/data"
)

func TestServiceImpl(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "API ServiceImpl unit Test Suite")
}

func createService(chatService chatserver.Service) Service {
	return NewServiceImpl(chatService)
}

var _ = ginkgo.Describe("ServiceImpl", func() {

	ginkgo.Context("GetMessages", func() {

		ginkgo.It("Return all messages", func() {

			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetMessages").Return(messages)
			responseMessages := service.GetMessages(9223372036854775807, 9223372036854775807)
			gomega.Expect(len(responseMessages)).To(gomega.Equal(3))
		})

		ginkgo.It("Return messages of a particular a user id", func() {

			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetMessages").Return(messages)
			responseMessages := service.GetMessages(1, 9223372036854775807)
			gomega.Expect(len(responseMessages)).To(gomega.Equal(1))
		})

		ginkgo.It("Return messages of a particular a room id", func() {

			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetMessages").Return(messages)
			responseMessages := service.GetMessages(9223372036854775807, 0)
			gomega.Expect(len(responseMessages)).To(gomega.Equal(3))
		})

		ginkgo.It("Return messages of a particular user in a room", func() {

			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetMessages").Return(messages)
			responseMessages := service.GetMessages(1, 0)
			gomega.Expect(len(responseMessages)).To(gomega.Equal(1))
		})
	})

	ginkgo.Context("PostMessage", func() {

		ginkgo.It("Post message returns success", func() {
			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetUser").Return(messages[1])
			chatServiceMock.On("Publish").Return(messages[1])
			responseMessage,_ := service.PostMessage(messages[1])
			gomega.Expect(responseMessage.ID).To(gomega.Equal(1))
		})

		ginkgo.It("Post message returns failure when roomId is not valid", func() {
			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetUser").Return(messages[1])
			chatServiceMock.On("Publish").Return(messages[1])
			_,err := service.PostMessage(data.Message{
				UserID: 1,
				RoomID: 3,
			})
			gomega.Expect(err.Error()).To(gomega.Equal("Room not found"))
		})

		ginkgo.It("Post message returns failure when userId is not valid", func() {
			chatServiceMock := &chatserver.ServiceMock{}
			service := createService(chatServiceMock)
			chatServiceMock.On("GetUser").Return(messages[1])
			chatServiceMock.On("Publish").Return(messages[1])
			_,err := service.PostMessage(data.Message{
				UserID: 5,
				RoomID: 0,
			})
			gomega.Expect(err.Error()).To(gomega.Equal("User not found"))
		})
	})
})

var users = []data.User {
	{
		ID: 0,
		Name: "System",
		ActiveRoom: 0,
	},
	{
		ID: 1,
		Name: "harish",
		ActiveRoom: 0,
	},
	{
		ID: 2,
		Name: "abhilash",
		ActiveRoom: 0,
	},
	{
		ID: 3,
		Name: "anusha",
		ActiveRoom: 0,
	},
}
var messages = []data.Message {
	{
		ID: 0,
		UserID: 1,
		RoomID: 0,
		UserName: "Rob",
		RoomName: "Default",
		Text: "hello",
		TimeStamp: "20190608172307",
	},
	{
		ID: 1,
		UserID: 2,
		RoomID: 0,
		UserName: "Bob",
		RoomName: "Default",
		Text: "hi this is abhilash",
		TimeStamp: "20190608172338",
	},
	{
		ID: 2,
		UserID: 3,
		RoomID: 0,
		UserName: "John",
		RoomName: "Default",
		Text: "hi this is anusha",
		TimeStamp: "20190608172359",
	},
}