package connections

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"chatServer/src/chatserver"
	"chatServer/src/chatserver/data"
	"chatServer/src/config"
)

// ServiceImpl struct for connections service
type ServiceImpl struct {
	chatService chatserver.Service
	config      *config.Config
}

// NewServiceImpl returns ServiceImpl
func NewServiceImpl(chatService chatserver.Service, config *config.Config) *ServiceImpl{
	return &ServiceImpl{
		chatService: chatService,
		config:	config,
	}
}

// HandleConnections handles the incoming connections
func (service *ServiceImpl) HandleConnections() {
	// listen for incoming tcp connections
	ln, err := net.Listen(service.config.ConnectionType, service.config.Host+":"+service.config.Port)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	defer ln.Close()

	// handle the incoming connections
	for {
		conn, err := ln.Accept()
		if err !=nil {
			log.Println("Error handling connection:", err.Error())
		}

		go service.handleConnection(conn)
	}
}


func (service *ServiceImpl) handleConnection(conn net.Conn) {
	log.Println("A new client joined")
	defer conn.Close()

	io.WriteString(conn, "Enter your username: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()

	user := service.chatService.CreateUser(scanner.Text())
	service.showCommands(conn)

	// handle writing back to connection
	go service.handleWriteToConnection(user, conn)

	// handle messages from the client
	for  {
		ok := scanner.Scan()
		message := scanner.Text()
		message = strings.TrimSpace(message)

		if len(message) > 0 {
			if message[0] == '/' { // check if it is a command
				service.handleCommands(message, conn, user)
			} else { // handle messages
				userStruct, _ := service.chatService.GetUser(user.ID)
				service.chatService.Publish(data.Input{
					Text: message,
					Room: userStruct.ActiveRoom,
				}, user.ID, false)
			}

			if !ok {
				break;
			}
		}
	}
}

// handleWriteToConnection writes back to connection
func (service *ServiceImpl) handleWriteToConnection(user data.User, conn net.Conn) {
	for {
		select {
			case msg := <- user.Output:
				io.WriteString(conn, msg)
			case <- user.Close:
				return
		}
	}
}

// handleCommands handles the commands by user and performs necessary actions
func (service *ServiceImpl) handleCommands(command string, conn net.Conn, user data.User) {
	switch {
	case command == "/help":
		service.showCommands(conn)
	case command == "/rooms":
		service.chatService.ListRooms(user.ID)
	case strings.HasPrefix(command, "/createroom"):
		if isCommandValid(command) {
			service.chatService.CreateRoom(strings.Replace(command, "/createroom ", "", -1), user.ID, user.Name)
		} else {
			sendOptionsMissingInfo(conn)
		}
	case strings.HasPrefix(command, "/subscribe"):
		if isCommandValid(command) {
			roomID, _ := strconv.Atoi(strings.Replace(command, "/subscribe ", "", -1))
			service.chatService.Subscribe(user.ID, roomID)
		} else {
			sendOptionsMissingInfo(conn)
		}
	case strings.HasPrefix(command, "/unsubscribe"):
		if isCommandValid(command) {
			roomID, _ := strconv.Atoi(strings.Replace(command, "/unsubscribe ", "", -1))
			service.chatService.UnSubscribe(user.ID, roomID)
		} else {
			sendOptionsMissingInfo(conn)
		}
	case strings.HasPrefix(command, "/switch"):
		if isCommandValid(command) {
			roomID, _ := strconv.Atoi(strings.Replace(command, "/switch ", "", -1))
			service.chatService.SwitchRoom(user.ID, roomID)
		} else {
			sendOptionsMissingInfo(conn)
		}
	case strings.HasPrefix(command, "/activeroom"):
		service.chatService.GetActiveRoom(user.ID)
	case command == "/quit":
		service.chatService.RemoveUser(user.ID)
		conn.Close()
		user.Close <- struct{} {}
	default:
		io.WriteString(conn, "Unknown Command!!!\n")
	}
}

// isCommandValid checks if command is valid or not
func isCommandValid(command string) bool {
	if len(strings.Split(command, " ")) != 2 {
		return false
	}
	return true
}

// sendOptionsMissingInfo checks if the options are missing or not in a command
func sendOptionsMissingInfo(conn net.Conn) {
	io.WriteString(conn, "Options missing!!!\n")
}

// showCommands shows the commands that are available to the user
func (service *ServiceImpl) showCommands(conn net.Conn) {
	commands :=
		`***Available commands***
/help - lists all the available commands
/rooms - lists all the available rooms
/createroom - creates new room - Ex: /createroom roomName
/subscribe - subscribes to a room - Ex: /subscribe roomId
/unsubscribe - unsubscribes from a room - Ex: /unsubscribe roomId
/switch - switches to a room - Ex: /switch roomId
/activeroom - displays the active room of a user - Ex: /activeroom
/quit` + "\n"
	io.WriteString(conn, commands)
}
