package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "encoding/json"

  "msg"
  "cliID"
)

func recData(c *net.TCPConn) {
  for {
    data := make([]byte, 4096)
    num, err := c.Read(data)
    if err != nil { panic(err) }

    if num > 0 {
      fmt.Println(string(data[:num]))
    }
  }
}

func sendRegularMesage(c *net.TCPConn, content string, author *cliID.CliID) {
  byt, err := json.Marshal( msg.NewMessage(0, content, author) )
  if err != nil { panic(err) }

  (*c).Write(byt)
}

func main() {
  tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil { panic(err) }
  defer conn.Close()

  go recData(conn)

  myID := cliID.CliID{
    IDnum: 0,
    Username: "egg",
  }

  scanner := bufio.NewScanner(os.Stdin)

  for {
    scanner.Scan()
    msg := scanner.Text()
    sendRegularMesage(conn, msg, &myID)

    if err != nil { panic(err) }
  }
}
