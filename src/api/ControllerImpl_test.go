package api

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"chatServer/src/chatserver/data"
)

func TestControllerImpl(t *testing.T) {

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "API ControllerImpl unit Test Suite")
}

func createController(service Service) Controller {
	return NewControllerImpl(service)
}

var _ = ginkgo.Describe("ControllerImpl", func() {

	const routeName = "/test"

	ginkgo.Context("GetMessages", func() {
		ginkgo.It("should get all the messages", func() {

			apiServiceMock := &ServiceMock{}
			controller := createController(apiServiceMock)

			url := routeName + "/rest/v1/messages"
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", url, nil)

			apiServiceMock.On("GetMessages",
				9223372036854775807, 9223372036854775807).Return(messages)
			controller.GetMessages(w, r)
			gomega.Expect(w.Code).To(gomega.Equal(200))
		})
	})

	ginkgo.Context("PostMessage", func() {
		ginkgo.It("should return bad request error when request body is not properly constructed", func() {
			apiServiceMock := &ServiceMock{}
			controller := createController(apiServiceMock)

			url := routeName + "/rest/v1/messages"
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", url, bytes.NewReader([]byte(`{}`)))

			apiServiceMock.On("PostMessage", data.Message{UserID: 0}).Return(nil, nil)
			controller.PostMessage(w, r)
			gomega.Expect(w.Code).To(gomega.Equal(400))
		})

		ginkgo.It("should call the controller and return 201", func() {
			apiServiceMock := &ServiceMock{}
			controller := createController(apiServiceMock)

			url := routeName + "/rest/v1/messages"
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", url, bytes.NewReader([]byte(`{"Text":"hello","userId": 1,"roomId": 0}`)))

			newMessage := data.Message{UserID: 1, Text: "hello", RoomID: 0}
			apiServiceMock.On("PostMessage", newMessage).Return(newMessage, nil)
			controller.PostMessage(w, r)
			gomega.Expect(w.Code).To(gomega.Equal(201))
		})

		ginkgo.It("should return 500 when userID or roomID does not exist", func() {
			apiServiceMock := &ServiceMock{}
			controller := createController(apiServiceMock)

			url := routeName + "/rest/v1/messages"
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", url, bytes.NewReader([]byte(`{"Text":"hello","userId": 5,"roomId": 0}`)))

			newMessage := data.Message{UserID: 5, Text: "hello", RoomID: 0}
			apiServiceMock.On("PostMessage", newMessage).Return(data.Message{}, errors.New(""))
			controller.PostMessage(w, r)
			gomega.Expect(w.Code).To(gomega.Equal(500))
		})
	})

})