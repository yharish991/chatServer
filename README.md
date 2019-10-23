# Chat Server
A simple chat server in Go where multiple clients can connect via telnet and send messages to the other connected clients.
Below are the features this chat server implements
- Client can send messages to other clients.
- Client can create a room in the chat server.
- Client can subscribe to a particular room.
- Client can unsubscribe to a particular room and the active room becomes Default room.
- Client can switch from one room to another to send messages to a particular room.
- Client can view the list of all rooms that are available.
- Messages sent by clients saved to local log file.
- REST APIs to post and query messages from chat server. 

## How it works?
- Chat server listens on a TCP port for the incoming TCP connections and handles those connections.
- Client establishes a TCP connection via telnet and sends the messages.
- Chat rooms can be shared between TCP clients.

## How to run the chat server
A Makefile has been created to make running the chat server easy. Below are the steps to run the chat server. Go version i used is `1.12.5`
- Go to /src folder in the project.
- Run the command `make deps`, this is to install the dependencies that are required for the chat server to run. This is a mandatory step.
- Run the command `make run` to run the chat server, this builds the project and runs it.

## How to connect as a client and use the chat server
Client can connect to the chat server by running the following command
`telnet 127.0.0.1 9080`

## Additional Makefile commands
Go to /src folder in the project.
- Use `make lint` to run golint on the Go files in the project.
- Use `make test` to run all the tests of this project.
- Use `make build` to compile packages and dependencies.
- Use `make all` to all the commands in the Makefile

## API Documentation
API Server runs on port 3000
### Post Messages API
Posts a message to a particular room for a particular user.
- ***URL***
`/rest/v1/messages`
- ***METHOD***
`POST`
- ***REQUEST BODY***
```$xslt
{
	"userId": 1,
	"roomId": 0,
	"text": "hello"
}
```
Required Body parameters
```$xslt
1. userId - int
2. text - string
```
- ***SUCCESSFUL RESPONSE***
```$xslt
{
    "id": 0,
    "userId": 1,
    "roomId": 0,
    "userName": "harish",
    "roomName": "Default",
    "text": "hello",
    "timestamp": "20190609115742"
}
```
- ***ERROR RESPONSE***
```$xslt
{
    "statusCode": 500,
    "message": "User not found"
}
```
- ***BAD REQUEST RESPONSE***
```$xslt
{
    "statusCode": 400,
    "message": "UserId or Text is empty"
}
```

### GET Messages API
API to query messages
- ***URL***
`/rest/v1/messages?userId={userId}&roomId={roomId}`
- ***METHOD***
`GET`
- ***QUERY PARAMETERS***
```$xslt
userId - returns the messages of a particular user - optional
roomId - returns the messages of a particular room - optional

when the query parameters are not specified, api returns all the messages
```
- ***SUCCESSFUL RESPONSE***
```$xslt
[
    {
        "id": 0,
        "userId": 0,
        "roomId": 1,
        "userName": "System",
        "roomName": "Tech",
        "text": "User John joined!!",
        "timestamp": "20190609121204"
    },
    {
        "id": 1,
        "userId": 1,
        "roomId": 1,
        "userName": "Bob",
        "roomName": "Default",
        "text": "Hello John!!",
        "timestamp": "20190609121221"
    },
    {
        "id": 2,
        "userId": 2,
        "roomId": 1,
        "userName": "John",
        "roomName": "Tech",
        "text": "Hey Bob!! whatsup?",
        "timestamp": "20190609121237"
    }
]
```

## Limitations/Constraints
- Right now as i don't persist the messages/users/rooms information to DB, i created an array of each of the objects to store the information.
- Id of each of the messages/users/rooms starts with 0 and gets incremented when a new message/user/room is created.
- userId - `0` is `System` and roomId - `0` is the `Default` room that gets created when the chat server is started.


