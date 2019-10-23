package chatserver

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"chatServer/src/chatserver/data"
)

// ServiceImpl struct for chat server service
type ServiceImpl struct {
	logFilePath string
	users []data.User
	rooms []data.Room
	messages []data.Message
	sync.RWMutex
}

// NewServiceImpl returns ServiceImpl
func NewServiceImpl(logFilePath string) *ServiceImpl {
	return &ServiceImpl{
		logFilePath: logFilePath,
	}
}


// Run preps and starts the chat server
func (service *ServiceImpl) Run() {
	service.createDefaultRoom()
	service.CreateUser("System") // System user
}

// CreateUser creates a new user
func (service *ServiceImpl) CreateUser(name string) data.User {
	service.Lock()
	defer service.Unlock()
	id := len(service.users)
	newUser := data.User{
		ID: id,
		Name: name,
		Output: make(chan string, 100),
	}
	newUser.ActiveRoom = service.rooms[0].ID // make the active room as Default room when user is created
	if (service.rooms[0].Users == nil) {
		service.rooms[0].Users = make(map[int]string)
	}
	service.rooms[0].Users[id] = name // add the created user to the Default room
	service.users = append(service.users, newUser)
	return newUser
}


// createDefaultRoom creates a new default room in chat chatserver
func (service *ServiceImpl) createDefaultRoom() {
	id := len(service.rooms)
	defaultRoom := data.Room{
		ID: id,
		Name: "Default",
		Users: make(map[int]string),
	}
	service.rooms = append(service.rooms, defaultRoom)
	log.Println("Default room created!!")
}


// Publish broadcasts the message to the users in the room
func (service *ServiceImpl) Publish(input data.Input, userID int, sysMessage bool) data.Message {
	service.Lock()
	defer service.Unlock()
	roomID := input.Room
	userList := service.rooms[roomID].Users

	timeStamp := service.getTimeStamp()
	formattedMessage := service.formatMessage(input,
		userID,
		service.users[userID].Name,
		service.rooms[input.Room].Name,sysMessage,
		timeStamp)

	// publish the message
	for id := range userList {
		userStruct := service.users[id]
		if id != userID && id != 0  && userStruct.Dead == false { // dont write message from self, to the system user and to dead user
			select {
				case userStruct.Output <- formattedMessage:
				case <-time.After(1 * time.Second):
					log.Printf("timeout sending to user %d", id)
			}
		}
	}
	service.logMessageToFile(formattedMessage)
	var uID int
	var uName string
	if sysMessage {
		uID = service.users[0].ID
		uName = service.users[0].Name
	} else {
		uID = service.users[userID].ID
		uName = service.users[userID].Name
	}
	savedMessage := service.saveMessage(uID, roomID, uName, service.rooms[roomID].Name, input.Text, timeStamp)
	return savedMessage
}


// Subscribe lets the user subscribe to a particular room
func (service *ServiceImpl) Subscribe(userID int, roomID int) {
	service.Lock()
	defer service.Unlock()
	// check if room is valid or not
	if roomID < len(service.rooms) {
		if service.rooms[roomID].Users[userID] == service.users[userID].Name { // check if already subscribed
			service.sendInfo("Already subscribed to room " + service.rooms[roomID].Name + "!!\n", userID)
		} else {
			service.rooms[roomID].Users[userID] = service.users[userID].Name
			service.sendInfo("Subscribed to " + service.rooms[roomID].Name + "!!\n", userID)
		}
	} else {
		service.sendInfo("Room " + strconv.Itoa(roomID) + " not found!!\n", userID)
	}
}


// UnSubscribe lets the user unsubscribe to a particular room
func (service *ServiceImpl) UnSubscribe(userID int, roomID int) {
	service.Lock()
	defer service.Unlock()
	// check if room is valid or not
	if roomID < len(service.rooms) {
		if service.rooms[roomID].Users[userID] == "" {
			service.sendInfo("User is not subscribed to " + service.rooms[roomID].Name + "!!\n", userID)
		} else {
			userList := service.rooms[roomID].Users
			delete(userList, userID)
			service.rooms[roomID].Users = userList
			if roomID == service.users[userID].ActiveRoom { // change the active room to Default if the user unsubscribes an active room
				service.users[userID].ActiveRoom = 0
			}
			service.sendInfo("Unsubscribed " + service.rooms[roomID].Name + "!!\n", userID)
		}
	} else {
		service.sendInfo("Room " + strconv.Itoa(roomID) + " not found!!\n", userID)
	}
}


// SwitchRoom lets the user switch to a particular room
func (service *ServiceImpl) SwitchRoom(userID int, roomID int) {
	service.Lock()
	defer service.Unlock()
	if roomID < len(service.rooms) {
		if service.users[userID].ActiveRoom == roomID {
			service.sendInfo("Already in room " + service.rooms[roomID].Name + "!!\n", userID)
		} else if service.rooms[roomID].Users[userID] != "" { // check if the user is subscribed to the room or not
			service.users[userID].ActiveRoom = roomID
			service.sendInfo("Switched to " + service.rooms[roomID].Name + "!!\n", userID)
		} else {
			service.sendInfo("Subscribe to " + service.rooms[roomID].Name + " before switching!!\n", userID)
		}
	} else {
		service.sendInfo("Room " + strconv.Itoa(roomID) + " not found!!\n", userID)
	}
}

// GetActiveRoom gets the active room of the user
func (service *ServiceImpl) GetActiveRoom(userID int) {
	service.RLock()
	defer service.RUnlock()
	activeRoomID := service.users[userID].ActiveRoom
	info := "Active room is " + service.rooms[activeRoomID].Name + " - " + strconv.Itoa(activeRoomID) + "!!\n"
	service.sendInfo(info, userID)
}


// ListRooms lists all the rooms in the chat chatserver
func (service *ServiceImpl) ListRooms(userID int) {
	service.RLock()
	defer service.RUnlock()
	var info string
	info = "List of rooms: \n"
	for _,room := range service.rooms {
		info = info + strconv.Itoa(room.ID) + "-" + room.Name + "\n"
	}
	service.sendInfo(info, userID)
}


// CreateRoom creates a new room in the chat server
func (service *ServiceImpl) CreateRoom(roomName string, userID int, userName string) {
	service.Lock()
	defer service.Unlock()
	// check if the room already exists
	for i := 0; i < len(service.rooms); i++ {
		if service.rooms[i].Name == roomName {
			service.sendInfo("Room with similar name already exists!!\n", userID)
			return
		}
	}
	room := data.Room{
		ID: len(service.rooms),
		Name: roomName,
	}
	if room.Users == nil {
		room.Users = make(map[int]string)
	}
	room.Users[userID] = userName
	service.rooms = append(service.rooms, room)
	service.sendInfo("Room " + roomName + " created!!\n", userID)
}


// GetUser gets a particular user details
func (service *ServiceImpl) GetUser(userID int) (data.User, bool) {
	service.RLock()
	defer service.RUnlock()
	// check if userID is valid or not
	if userID < len(service.users) {
		return service.users[userID], true
	}
	return data.User{}, false
}


// GetRoom gets a particular room details
func (service *ServiceImpl) GetRoom(roomID int) (data.Room, bool) {
	service.RLock()
	defer service.RUnlock()
	// check if roomID is valid or not
	if roomID < len(service.rooms) {
		return service.rooms[roomID], true
	}
	return data.Room{}, false
}

// GetMessages returns all the messages
func (service *ServiceImpl) GetMessages() []data.Message {
	service.RLock()
	defer service.RUnlock()
	return service.messages
}

// GetUsers returns all the users
func (service *ServiceImpl) GetUsers() []data.User {
	service.RLock()
	defer service.RUnlock()
	return service.users
}

// GetRooms returns all the rooms
func (service *ServiceImpl) GetRooms() []data.Room {
	service.RLock()
	defer service.RUnlock()
	return service.rooms
}

// RemoveUser marks the user as dead
func (service *ServiceImpl) RemoveUser(userID int) {
	service.RLock()
	defer service.RUnlock()
	service.users[userID].Dead = true;
}

// formatMessage formats the message to a particular format
func (service *ServiceImpl) formatMessage(
	input data.Input,
	userID int,
	userName string,
	roomName string,
	sysMessage bool,
	timeStamp string) string {
	var user string
	if sysMessage {
		user = "System"
	} else {
		user = userName
	}
	return fmt.Sprintf("%s %s |%s| %s\n",
		timeStamp, "Room:" + roomName, user, input.Text)
}


// getTimeStamp gets the timestamp
func (service *ServiceImpl) getTimeStamp() string {
	t := time.Now()
	return t.Format("20060102150405")
}


// logMessageToFile logs the message to the log file
func (service *ServiceImpl) logMessageToFile(message string) {
	file, err := os.OpenFile(
		path.Join(service.logFilePath),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	if _, err := file.WriteString(message); err != nil {
		panic(err)
	}
}


// saveMessage saves the message
func (service *ServiceImpl) saveMessage(userID int, roomID int, userName string, roomName string, text string, timeStamp string)  data.Message{
	newMessage := data.Message{
		ID: len(service.messages),
		UserID: userID,
		RoomID: roomID,
		UserName: userName,
		RoomName: roomName,
		Text: text,
		TimeStamp: timeStamp,
	}
	service.messages = append(service.messages, newMessage)
	return newMessage
}


// sendInfo sends the info to a particular user
func (service *ServiceImpl) sendInfo(info string, userID int) {
	user := service.users[userID]
	user.Output <- info
}