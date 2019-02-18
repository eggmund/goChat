package main

import (
  "net"
  "fmt"
  "io"
)

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
      fmt.Println(c.RemoteAddr(), ":", string(buff[:num]))
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
    fmt.Println("Connected to:", ln.Addr())

    go handleConnection(conn)
  }
}
