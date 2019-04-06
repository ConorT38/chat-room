package main

import (
  "fmt"
  "net"
  "os"
  "regexp"
  "time"
)

const (
  CONN_HOST      = "localhost"
  CONN_PORT      = "80"
  CONN_TYPE      = "tcp"
  ERROR          = "Error"
  INFO           = "Info"
  CREATE_SUCCESS = "successfully created room"
  JOIN_SUCCESS   = "successfully joined room"
)

type Message struct {
  from User
  room string
  time string
}

type User struct {
  name string
  room string
  conn net.Conn
}

var connections []net.Conn
var rooms = make(map[string][]net.Conn)
func main() {
  l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

  if err != nil {
    log(err.Error(), ERROR)
    os.Exit(1)
  }

  defer l.Close()
  log("Listening on "+CONN_HOST+":"+CONN_PORT, INFO)

  for {

    conn, err := l.Accept()
    connections = append(connections, conn)

    if err != nil {
      log(err.Error(), ERROR)
      handleDisconnection()
      os.Exit(1)
    }

    go handleRequest(conn)

  }
}

func broadcast(message string, conn net.Conn) {
  room := getFromBrackets(message)
  for _, connection := range rooms[room] {
    if conn != connection {
      sendMessage(connection, message)
    }
  }
}

func handleRequest(conn net.Conn) {
  buf := make([]byte, 1024)

  for {
    _, err := conn.Read(buf)
    message := string(buf)

    if err != nil {
      log(err.Error(), ERROR)
      break
    }

    log(message, INFO)
    if message[:11] == "[JOIN_ROOM]" {
      room := getFromBrackets(message[11:])
      if joinRoom(conn, room) {
        sendMessage(conn, JOIN_SUCCESS)
      } else if createRoom(conn, room) {
        sendMessage(conn, CREATE_SUCCESS)
      }
    } else {
      go broadcast(message, conn)
    }
  }
  conn.Close()
}

func createRoom(conn net.Conn, room string) bool {
  if _, ok := rooms[room]; ok {
    sendMessage(conn, "[SERVER]: Room exists")
    return false
  } else {
    log("room created:"+room, INFO)
    rooms[room] = append(rooms[room], conn)
    return true
  }
}

func joinRoom(conn net.Conn, room string) bool {
  if _, ok := rooms[room]; ok {
    rooms[room] = append(rooms[room], conn)
    go broadcast("[SERVER]: New connection", conn)
    return true
  } else {
    sendMessage(conn, "[SERVER]: Couldn't join room")
    return false
  }
}

func sendMessage(conn net.Conn, message string) {
  conn.Write([]byte(message + "\n"))
  log("sent message "+message, INFO)
}

func log(message string, type_ string) {
  log_time := "[" + time.Now().Format("3:04PM") + "]"
  ERROR := "[ERROR]"
  INFO := "[INFO]"

  if type_ != "Error" {
    fmt.Println(INFO + log_time + ": " + message)
  } else {
    fmt.Println(ERROR + log_time + ": " + message)
  }
}

func getFromBrackets(str string) string {
  log(str, INFO)
  r, _ := regexp.Compile("^(.*?)]")
  match := r.FindString(str)
  match = match[1:]
  return match[:len(match)-1]
}

func handleDisconnection() {
  log("User disconnected", INFO)
}
