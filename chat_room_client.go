package main

import( 
 "net"
 "fmt"
 "bufio"
 "os"
 "time"
 "os/signal"
  "syscall"
 
 )

const (
  TIME_FORMAT = "3:04PM"
  ENTER_ALIAS = "Enter alias here: "
  ENTER_MESSAGE = "Enter in a message below:\n"
  CONN_TYPE = "tcp"
  CONN_HOST = "localhost"
  CONN_PORT = "80"
)

type Message struct{
  from User
  room string
  time string
}

type User struct{
  name string
  room string
  conn net.Conn
}

func main() {

  // connect to this socket
  conn, _ := net.Dial(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
  user := login(conn)

  for { 
    handleQuit(user)
    // read in input from stdin
    inputLabel(user)
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    time_sent := "["+time.Now().Format(TIME_FORMAT)+"]"
    
    // send to socket
    if text != "" {
    sendMessage(user, user.room+time_sent+user.name+text + "\n")
  }
    
  }
}

func listen(user User){
  for{
    message, _ := bufio.NewReader(user.conn).ReadString('\n')
    fmt.Println("\n"+string(message))
    inputLabel(user)
  }
}

func login(conn net.Conn) User{
  var name string 
  var room string

  InputLoop:
    for name == "" || room == ""{ 
      fmt.Print(ENTER_ALIAS) 
      fmt.Scanln(&name)
      if name == ""{
        log("Name can't be blank")
        continue InputLoop
      }
      name = "["+name+"]: "

      fmt.Print("Create or Join a room: ")
      fmt.Scanln(&room)
      if room == "" {
        log("room can't be blank")
        continue InputLoop
      }
      room = "["+room+"]"
  }

  user := User{name: name, room: room, conn: conn} 
  joinRoom(user)
  
  go listen(user) 

  fmt.Print(ENTER_MESSAGE)

  return user
}

//TODO: add user object instead of conn
func sendMessage(user User, message string){
  fmt.Fprintf(user.conn, message)
}

func joinRoom(user User) {
  request := "[JOIN_ROOM]"+user.room
  sendMessage(user,request)
}

func log(message string){
  fmt.Println(message+"\n")
}

func handleQuit(user User){
   c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        sendMessage(user,"[SERVER]: "+user.name[:len(user.name)-2]+" disconnected")
        log("You disconnected from "+user.room)
        os.Exit(1)
    }()
}

func inputLabel(user User){
  f := bufio.NewWriter(os.Stdout)
  defer f.Flush()
  f.Write([]byte(user.room+user.name))
}
