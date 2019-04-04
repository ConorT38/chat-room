package main

import (
    "fmt"
    "net"
    "os"
    "time"
    "regexp"
    "strings"
    )

const (
    CONN_HOST = "localhost"
    CONN_PORT = "80"
    CONN_TYPE = "tcp"
    ERROR = "Error"
    INFO = "Info"
    SUCCESS = "success"
)

var connections []net.Conn
var rooms = make(map[string][]net.Conn)

func main() {
    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    
    if err != nil {
        log(err.Error(),ERROR)
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    log("Listening on " + CONN_HOST + ":" + CONN_PORT,INFO)

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        connections = append(connections,conn)

        if err != nil {
            log(err.Error(),ERROR)
            os.Exit(1)
        }

        // Handle connections in a new goroutine.
        go handleRequest(conn)
       
    }
}

// Handles incoming requests.
func broadcast(message string, conn net.Conn) {
  
  r, _ := regexp.Compile("^(.*?)]")
  match := r.FindString(message)
  match = match[1:]
  match = match[:len(match)-1]

  for k, v := range rooms { 
    fmt.Printf("key[%s] value[%s]\n", k, v)
}
  for _, connection := range rooms[match]{
    if conn != connection {
    connection.Write([]byte(message))
  }
}
}

func handleRequest(conn net.Conn){
    // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  for{
  _, err := conn.Read(buf)
  message := string(buf)

  if err != nil {
    log(err.Error(),ERROR)
    break
  }

  log(message,INFO)
  if message[:11] == "[JOIN_ROOM]"{ 
    //TODO: grab only characters and not whitespace from room
    r, _ := regexp.Compile("(\\S)*")
    match := r.FindString(message)
    log("["+match+"]",INFO)
    if joinRoom(conn,message[11:]){
     sendMessage(conn,SUCCESS)
    }else if createRoom(conn,message[11:]) && joinRoom(conn,message[11:]){
      sendMessage(conn,SUCCESS)
    }
  } else{
    go broadcast(message,conn)
}
  // Close the connection when you're done with it.
  }
  conn.Close()
}

func createRoom(conn net.Conn, room string) bool{
  fmt.Printf("["+room+"]")
  if _, ok := rooms[room]; ok {
      sendMessage(conn,"failed")
      return false
  } else{
    log("room created:"+ room,INFO)
    rooms[room] = append(rooms[room],conn)
    return true
  }
}

func joinRoom(conn net.Conn, room string) bool{
  room = strings.TrimSpace(room[:len(room)-1])
  if _, ok := rooms[room]; ok {
    rooms[room] = append(rooms[room],conn)
    log("Joined room "+room,INFO)
    return true
  } else{
    sendMessage(conn, "failed")
    return false
  }
}

func sendMessage(conn net.Conn, message string){
  conn.Write([]byte(message+"\n"))
  log("sent message "+message,INFO)
}

func log(message string, type_ string){
      log_time := "["+time.Now().Format("3:04PM")+"]"
      ERROR := "[ERROR]"
      INFO := "[INFO]"

      if type_ != "Error" {
        fmt.Println(INFO+log_time+": "+message)
      }else{
        fmt.Println(ERROR+log_time+": "+message)
      }
}