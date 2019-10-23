package chatserver

import (
	"chatServer/src/chatserver/data"
	"path"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"chatServer/testhelpers"
)

func TestServiceImpl(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Chat chatserver ServiceImpl unit Test Suite")
}

func createService(logFilePath string) Service {
	return NewServiceImpl(logFilePath)
}

var _ = ginkgo.Describe("ServiceImpl", func() {

	ginkgo.Context("Run", func() {

		ginkgo.It("Creates a default room and system user", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			gomega.Expect(len(service.GetUsers())).To(gomega.Equal(1))
			gomega.Expect(len(service.GetUsers())).To(gomega.Equal(1))
			gomega.Expect(service.GetUsers()[0].Name).To(gomega.Equal("System"))
			gomega.Expect(service.GetRooms()[0].Name).To(gomega.Equal("Default"))
		})
	})

	ginkgo.Context("CreateUser", func() {

		ginkgo.It("Creates a new user and returns the user object", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			user := service.CreateUser("TestUser")
			gomega.Expect(user.Name).To(gomega.Equal("TestUser"))
		})
	})

	ginkgo.Context("CreateRoom", func() {

		ginkgo.It("Creates a new room", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			user, _ := service.GetUser(1)
			gomega.Expect(len(service.GetRooms())).To(gomega.Equal(2))
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
		})

		ginkgo.It("Cannot create a room with same name", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			user, _ := service.GetUser(1)
			gomega.Expect(len(service.GetRooms())).To(gomega.Equal(2))
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(<-user.Output).To(gomega.Equal("Room with similar name already exists!!\n"))
		})
	})

	ginkgo.Context("UnSubscribe", func() {

		ginkgo.It("UnSubscribes to a room", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			service.UnSubscribe(1, 1)
			user, _ := service.GetUser(1)
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(<-user.Output).To(gomega.Equal("Unsubscribed Tech!!\n"))
		})
	})

	ginkgo.Context("Subscribe", func() {

		ginkgo.It("Subscribes to a room", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			user, _ := service.GetUser(1)
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(service.GetRooms()[1].Users[1]).To(gomega.Equal("TestUser"))
		})

		ginkgo.It("sends info to the user that the user is already subscribed when subscribing to a same room", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			service.Subscribe(1, 1)
			user, _ := service.GetUser(1)
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(<-user.Output).To(gomega.Equal("Already subscribed to room Tech!!\n"))
		})
	})

	ginkgo.Context("SwitchRoom", func() {

		ginkgo.It("switches to a room", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			service.SwitchRoom(1, 1)
			user, _ := service.GetUser(1)
			gomega.Expect(service.GetUsers()[1].ActiveRoom).To(gomega.Equal(1))
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(<-user.Output).To(gomega.Equal("Switched to Tech!!\n"))
		})
	})

	ginkgo.Context("ListRooms", func() {

		ginkgo.It("lists all the rooms", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			service.ListRooms(1)
			user, _ := service.GetUser(1)
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(<-user.Output).To(gomega.ContainSubstring("Tech"))
		})
	})

	ginkgo.Context("GetUser", func() {

		ginkgo.It("Gets the user details", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			user, ok := service.GetUser(1)
			gomega.Expect(user.Name).To(gomega.ContainSubstring("TestUser"))
			gomega.Expect(ok).To(gomega.Equal(true))
		})
	})

	ginkgo.Context("GetRoom", func() {

		ginkgo.It("Gets the room details", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			room,_ := service.GetRoom(1)
			user, _ := service.GetUser(1)
			gomega.Expect(<-user.Output).To(gomega.Equal("Room Tech created!!\n"))
			gomega.Expect(room.Name).To(gomega.ContainSubstring("Tech"))
		})
	})

	ginkgo.Context("Publish", func() {

		ginkgo.It("Publishes the message to all users in the room", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			newUser := service.CreateUser("Bob")
			service.Publish(data.Input{
				0,
				"Hello!!",
			}, 1, false)
			user, _ := service.GetUser(newUser.ID)
			gomega.Expect(<-user.Output).To(gomega.ContainSubstring("Hello!!"))
		})
	})

	ginkgo.Context("GetMessages", func() {

		ginkgo.It("Gets all the messages", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateUser("Bob")
			service.Publish(data.Input{
				0,
				"Hello!!",
			}, 1, false)
			gomega.Expect(len(service.GetMessages())).To(gomega.Equal(1))
		})
	})

	ginkgo.Context("GetUsers", func() {

		ginkgo.It("Gets all the users", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateUser("Bob")
			gomega.Expect(len(service.GetUsers())).To(gomega.Equal(3))
		})
	})

	ginkgo.Context("GetRooms", func() {

		ginkgo.It("Gets all the rooms", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.CreateRoom("Tech", 1, "TestUser")
			gomega.Expect(len(service.GetRooms())).To(gomega.Equal(2))
		})
	})

	ginkgo.Context("RemoveUser", func() {

		ginkgo.It("marks a particular user as dead", func() {
			service := createService(path.Join(testhelpers.GetServerRootDir(), "/logs/messages.log"))
			service.Run()
			service.CreateUser("TestUser")
			service.RemoveUser(1)
			user, _ := service.GetUser(1)
			gomega.Expect(len(service.GetUsers())).To(gomega.Equal(2))
			gomega.Expect(user.Dead).To(gomega.Equal(true))
		})
	})
})
