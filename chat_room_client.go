package main

import( 
 "net"
 "fmt"
 "bufio"
 "os"
 "time"
 
 )

const (
  TIME_FORMAT = "3:04PM"
  ENTER_ALIAS = "Enter alias here: "
  ENTER_MESSAGE = "Enter in a message below:\n"
  CONN_TYPE = "tcp"
  CONN_HOST = "localhost"
  CONN_PORT = "80"
)

var name string 
var room string

func main() {

  // connect to this socket
  conn, _ := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
  login(conn)
  room = "["+room+"]"
  
  for { 
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    time_sent := "["+time.Now().Format(TIME_FORMAT)+"]"
    
    // send to socket
    if text != "" {
    sendMessage(conn, room+time_sent+name+text + "\n")
  }
    
  }
}

func listen(conn net.Conn){
  for{
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Println("\n"+string(message))
  }
}

func login(conn net.Conn){
  fmt.Print(ENTER_ALIAS) 
  fmt.Scanln(&name)
  name = "["+name+"]: "

  fmt.Print("Create or Join a room: ")
  fmt.Scanln(&room)
  joinRoom(conn,room)
  
  go listen(conn) 

  fmt.Print(ENTER_MESSAGE)
}

func sendMessage(conn net.Conn, message string){
  fmt.Fprintf(conn, message)
}

func joinRoom(conn net.Conn, room string) {
  request := "[JOIN_ROOM]" + room
  sendMessage(conn,request)
}
// func joinRoom(conn net.Conn, room string) bool{
//   request := "[JOIN_ROOM]" + room
//   var message string 

//   sendMessage(conn,request)
//   go func(conn net.Conn) {
//     for{
//       message, _ = bufio.NewReader(conn).ReadString('\n')
//       log(message)
//       if message == "success" {
//         log("You have succesfully joined room: "+room)
//         break
//       }
//     }
//     }(conn)

//     if message == "success" {
//       return true
//     }
//     log("dang")
//     return false
// }

// func createRoom(conn net.Conn, room string) bool{
//   request := "[CREATE_ROOM]" + room 
//   //sendMessage(conn,request)
//   var status string
//   fmt.Fprintf(conn, request)
//   //go func(conn net.Conn) {
//     for{
//       message, _, err := bufio.NewReader(conn).ReadString('\n')
//       if err !=nil{
//         log(err.Error())
//       }
//       if message == "success" {
//         log("You have succesfully joined room: "+room)
//         status = message
//         break
//       }
//     }
//     //}(conn)

//     if status == "success" {
//       return true
//     }
//     log("dang")
//     return false
// }

func log(message string){
  fmt.Println(message+"\n")
}