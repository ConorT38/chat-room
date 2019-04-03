package main

import "net"
import "fmt"
import "bufio"
import "os"

func main() {

  // connect to this socket
  conn, _ := net.Dial("tcp", "localhost:80")
  fmt.Print("Enter alias here: ")
  var name string 
  fmt.Scanln(&name)
  name = "["+name+"]: "
  
  go listen(conn)

  for { 
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(name)
    text, _ := reader.ReadString('\n')
    // send to socket
    fmt.Fprintf(conn, name+text + "\n")
    // listen for reply
    
  }
}

func listen(conn net.Conn){
  for{
    message, _ := bufio.NewReader(conn).ReadString('\n')
    fmt.Println("\n"+message)
  }
}