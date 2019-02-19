package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "encoding/json"

  "msg"
  "cliTools"
)

func recData(c *net.TCPConn) []byte {
  for {
    data := make([]byte, 4096)
    num, err := c.Read(data)
    if err != nil { panic(err) }

    if num > 0 {
      return data[:num]
    }
  }
}

func recNormMessages(c *net.TCPConn) {
  for {
    data := recData(c)
    var m msg.Message
    err := json.Unmarshal(data, &m)
    if err != nil { panic(err) }

    if m.Type == 0 {
      fmt.Printf("%s:: %s\n", m.Author.Username, m.Content)
    }
  }
}

func sendRegularMesage(c *net.TCPConn, content string, author *cliTools.CliID) {
  byt, err := json.Marshal( msg.NewMessage(0, content, author) )
  if err != nil { panic(err) }

  (*c).Write(byt)
}

func getMyIDnum(c *net.TCPConn, username string) int {
  byt, err := json.Marshal( msg.NewMessage(1, username, nil) )
  if err != nil { panic(err) }

  (*c).Write(byt)

  return 0
}

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  println("Enter username:")
  scanner.Scan()
  username := scanner.Text()

  tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil { panic(err) }
  defer conn.Close()

  println("Connected.")
  go recNormMessages(conn)

  //getMyIDnum(conn, username)

  myID := cliTools.CliID{
    IDnum: 0,
    Username: username,
  }


  for {
    scanner.Scan()
    msg := scanner.Text()
    sendRegularMesage(conn, msg, &myID)

    if err != nil { panic(err) }
  }
}
