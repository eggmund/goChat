package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
)

func recData(c net.Conn) {
  for {
    data := make([]byte, 0)
    _, err := c.Read(data)
    if err != nil { panic(err) }
    if len(data) != 0 {
      fmt.Println(data)
    }
  }
}

func main() {
  tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil { panic(err) }
  defer conn.Close()

  go recData(conn)

  scanner := bufio.NewScanner(os.Stdin)

  for {
    scanner.Scan()
    msg := scanner.Text()
    _, err = conn.Write([]byte(msg))
    if err != nil { panic(err) }

    fmt.Println(conn.LocalAddr(), ":", msg)
  }
}
