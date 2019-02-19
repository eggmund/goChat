package main

import (
  "net"
  "fmt"
  "io"

  "cliTools"
)

var (
  clients []*cliTools.CliData
  messages []string
)

func broadcast(clients []*cliTools.CliData, byt []byte) {
  for i := 0; i < len(clients); i++ {
    _, err := (*clients[i]).Conn.Write(byt)
    if err != nil && err != io.EOF {
      panic(err)
    }
  }
}

func initiateCli(c *net.Conn, username string) {
  
}

func recData(c *net.Conn) []byte {
  for {
    data := make([]byte, 4096)
    num, err := *c.Read(data)
    if err != nil { panic(err) }

    if num > 0 {
      return data[:num]
    }
  }
}

func handleConnection(c net.Conn) {
  initiateCli(&c)

  for c.RemoteAddr() != nil {
    data := recData(&c)
    var m msg.Message
    err := json.Unmarshal(data, &m)
    if err != nil { panic(err) }

    if m.Type == 0 {
      broadcast(clients, data)
    } else if m.Type == 1 {
      initiateCli(&c, m.Content)
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
    conn, err := ln.Accept()
    if err != nil { panic(err) }
    clients = append(clients, &conn)
    fmt.Println("Connected to:", ln.Addr())

    go handleConnection(conn)
  }
}
