package data

// User is a User Object
type User struct {
	ID            int
	Name          string
	ActiveRoom    int
	Output chan   string
	Close chan    struct{}
	Dead          bool
}

// Input is a Input Object
type Input struct {
	Room          int
	Text          string
}

// Room is a Room Object
type Room struct {
	ID            int
	Name          string
	Users         map[int]string
}

// Message is a Message Object
type Message struct {
	ID            int        `json:"id"`
	UserID        int        `json:"userId"`
	RoomID        int        `json:"roomId"`
	UserName      string     `json:"userName"`
	RoomName      string     `json:"roomName"`
	Text          string     `json:"text"`
	TimeStamp     string     `json:"timestamp"`
}
