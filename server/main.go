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
  var err error
  for m.Type != 1 {
    println("Waiting for message of username.")
    m, _, err = recMsg(c)
    if err != nil && err != io.EOF {
      panic(err)
    }
  }
  usrname, ok := m.Content.(string)
  if !ok { panic("Not ok :(") }
  println("Got username:", usrname)

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

func recData(c *net.TCPConn) ([]byte, error) {
  for {
    data := make([]byte, 4096)
    num, err := c.Read(data)
    if err != nil {
      return []byte{}, err
    }

    if num > 0 {
      return data[:num], nil
    }
  }
}

func recMsg(c *net.TCPConn) (msg.Message, []byte, error) {
  data, err := recData(c)
  var m msg.Message

  if err != nil {
    if err != io.EOF {
      panic(err)
    } else {
      return m, []byte{}, err
    }
  }
  err = json.Unmarshal(data, &m)
  if err != nil { panic(err) }

  return m, data, nil
}

func handleConnection(c *net.TCPConn) {
  cli := initiateCli(c)
  println("Initiated client")
  clients = append(clients, cli)
  fmt.Println("New client added:", cli.ID.Username)

  for {
    m, rawData, err := recMsg(c)
    if err == io.EOF {
      e := c.Close()
      if e != nil { panic(err) }
      removeCliAtID(cli.ID.IDnum)
      fmt.Println("Connection closed to:", c.RemoteAddr())
      break
    }
    if m.Type == 0 {
      broadcast(clients, rawData)
    }
  }
}

func removeCliAtID(id int) {
  var ind int = -1
  for i := 0; i < len(clients) && ind == -1; i++ {
    if clients[i].ID.IDnum == id {
      ind = i
    }
  }
  if ind != -1 {
    clients = append(clients[:ind], clients[ind+1:]...)
  }
}

func main() {
  addr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  ln, err := net.ListenTCP("tcp", addr)
  if err != nil { panic(err) }

  fmt.Println("Ready.")

  for {
    var conn *net.TCPConn
    conn, err = ln.AcceptTCP()
    if err != nil { panic(err) }
    fmt.Println("Connected to:", ln.Addr())

    go handleConnection(conn)
  }
}
