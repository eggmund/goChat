package main

import (
  "net"
  "fmt"
  "io"
  "encoding/json"

  "msg"
  "cliTools"
)

var (
  clients []*cliTools.CliData
  messages []string
)

func broadcast(clients []*cliTools.CliData, byt []byte) {
  fmt.Println("Broadcasting:", string(byt))

  for i := 0; i < len(clients); i++ {
    _, err := (*clients[i]).Conn.Write(byt)
    if err != nil && err != io.EOF {
      panic(err)
    }
  }
}

func initiateCli(c *net.TCPConn) *cliTools.CliData {
  m := *new(msg.Message)
  for m.Type != 1 {
    m, _ = recMsg(c)
  }
  usrname, ok := m.Content.(string)
  if !ok { panic("Not ok :(") }

  byt, err := json.Marshal( msg.NewMessage(1, len(clients), nil) )
  if err != nil { panic(err) }
  c.Write(byt)  // Send id num

  return &cliTools.CliData {
    Conn: c,
    ID: &cliTools.CliID {
      IDnum: len(clients),
      Username: usrname,
    },
  }
}

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

func recMsg(c *net.TCPConn) (msg.Message, []byte) {
  data := recData(c)
  var m msg.Message
  err := json.Unmarshal(data, &m)
  if err != nil { panic(err) }

  return m, data
}

func handleConnection(c *net.TCPConn) {
  clients = append(clients, initiateCli(c))
  fmt.Println("New client added:", clients[len(clients)-1].ID.Username)

  for c.RemoteAddr() != nil {
    m, rawData := recMsg(c)
    if m.Type == 0 {
      broadcast(clients, rawData)
    }
  }

  fmt.Println("Connection closed to:", c.RemoteAddr())
}

func main() {
  addr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  ln, err := net.ListenTCP("tcp", addr)
  if err != nil { panic(err) }

  for {
    var conn *net.TCPConn
    conn, err = ln.AcceptTCP()
    if err != nil { panic(err) }
    fmt.Println("Connected to:", ln.Addr())

    go handleConnection(conn)
  }
}
