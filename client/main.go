package main

import (
  "fmt"
  "net"
  "bufio"
  "os"

  "cliID"
  "json"
)

func recData(c net.Conn) {
  for {
    data := make([]byte, 4096)
    num, err := c.Read(data)
    if err != nil { panic(err) }

    if num > 0 {
      fmt.Println(string(data[:num]))
    }
  }
}

func getID(c *net.Conn) cliID.CliID {
  
}

func main() {
  tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil { panic(err) }
  defer conn.Close()

  getID(&conn)
  go recData(conn)

  scanner := bufio.NewScanner(os.Stdin)

  for {
    scanner.Scan()
    msg := scanner.Text()
    _, err = conn.Write([]byte(msg))
    if err != nil { panic(err) }
  }
}
