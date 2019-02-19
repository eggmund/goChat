package main

import (
  "net"
  "fmt"
  "io"
)

var (
  clients []*net.Conn
  messages []string
)

func broadcast(clients []*net.Conn, msg string) {
  for i := 0; i < len(clients); i++ {
    println("Broadcasting:", msg)
    _, err := (*clients[i]).Write([]byte(msg))
    if err != nil && err != io.EOF {
      panic(err)
    }
  }
}

func handleConnection(c net.Conn) {
  for c.RemoteAddr() != nil {
    buff := make([]byte, 4096)
    num, err := c.Read(buff)
    if err != nil {
      if err != io.EOF {
        panic(err)
      } else {
        break
      }
    }

    if num > 0 {
      msg := fmt.Sprintf("%v:: %s", c.RemoteAddr(), string(buff[:num]))
      broadcast(clients, msg)
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
