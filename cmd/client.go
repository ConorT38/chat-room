package chat_room

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	TIME_FORMAT   = "3:04PM"
	ENTER_ALIAS   = "Enter alias here: "
	ENTER_MESSAGE = "Enter in a message below:\n"
	CONN_PROTOCOL = "tcp"
	HOST          = "localhost"
	PORT          = "80"
)

type ChatMessage struct {
	from ChatUser
	room string
	time string
}

type ChatUser struct {
	name string
	room string
	conn net.Conn
}

func main() {

	conn, _ := net.Dial(CONN_PROTOCOL, HOST+":"+PORT)
	user := login(conn)

	for {
		handleQuit(user)
		inputLabel(user)

		fmt.Print(user.room + user.name)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		time_sent := "[" + time.Now().Format(TIME_FORMAT) + "]"

		if text != "" {
			sendMessage(user, user.room+time_sent+user.name+text+"\n")
		}

	}
}

func listen(user ChatUser) {
	for {
		message, _ := bufio.NewReader(user.conn).ReadString('\n')
		fmt.Println("\n" + string(message))
		//inputLabel(user)
	}
}

func login(conn net.Conn) ChatUser {
	var name string
	var room string

InputLoop:
	for name == "" || room == "" {
		fmt.Print(ENTER_ALIAS)
		fmt.Scanln(&name)
		if name == "" {
			log("Name can't be blank")
			continue InputLoop
		}
		name = "[" + name + "]: "

		fmt.Print("Create or Join a room: ")
		fmt.Scanln(&room)
		if room == "" {
			log("room can't be blank")
			continue InputLoop
		}
		room = "[" + room + "]"
	}

	user := ChatUser{name: name, room: room, conn: conn}
	joinRoom(user)

	go listen(user)

	fmt.Print(ENTER_MESSAGE)

	return user
}

func sendMessage(user ChatUser, message string) {
	fmt.Fprintf(user.conn, message)
}

func joinRoom(user ChatUser) {
	request := "[JOIN_ROOM]" + user.room
	sendMessage(user, request)
}

func log(message string) {
	fmt.Println(message + "\n")
}

func handleQuit(user ChatUser) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		sendMessage(user, "[SERVER]: "+user.name[:len(user.name)-2]+" disconnected")
		log("You disconnected from " + user.room)
		os.Exit(1)
	}()
}

func inputLabel(user ChatUser) {
	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()
	f.Write([]byte(user.room + user.name))
}
