package main

import (
  "fmt"
  "net"
  "bufio"
  "os"
  "encoding/json"
  "github.com/gen2brain/raylib-go/raylib"

  "msg"
  "cliTools"
)

var (
  SCREEN_W int32 = 1000
  SCREEN_H int32 = 600

  myID *cliTools.CliID
)

func recData(c *net.TCPConn) []byte {
  var (
    num int
    data []byte
    err error
  )

  for num == 0 {
    data = make([]byte, 4096)
    num, err = c.Read(data)
    if err != nil { panic(err) }
  }

  return data[:num]
}

func recNormMessages(c *net.TCPConn) {
  for {
    data := recData(c)
    m := convDataToMsg(data)

    if m.Type == 0 {
      if ms, ok := m.Content.(string); ok {
        fmt.Printf("%s:: %s\n", m.Author.Username, ms)
      }
    }
  }
}

func convDataToMsg(data []byte) msg.Message {
  var m msg.Message
  err := json.Unmarshal(data, &m)
  if err != nil { panic(err) }

  return m
}

func sendRegularMesage(c *net.TCPConn, content string, author *cliTools.CliID) {
  byt, err := json.Marshal( msg.NewMessage(0, content, author) )
  if err != nil { panic(err) }

  (*c).Write(byt)
}

func getMyID(c *net.TCPConn, username string) *cliTools.CliID {
  byt, err := json.Marshal( msg.NewMessage(1, username, nil) )
  if err != nil { panic(err) }
  c.Write(byt)

  m := *new(msg.Message)
  for m.Type != 1 {
    println("Waiting for message.")
    m = convDataToMsg(recData(c))
  }
  idNum, ok := m.Content.(float64)
  if !ok { panic("Type assertion failed") }

  println("My id num:", int(idNum))

  return &cliTools.CliID {
    IDnum: int(idNum),
    Username: username,
  }
}

func initConnection(username string) *net.TCPConn {
  tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:6779")
  if err != nil { panic(err) }

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  if err != nil { panic(err) }
  println("Connected.")

  myID = getMyID(conn, username)
  go recNormMessages(conn)

  return conn
}

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  println("Enter username:")
  scanner.Scan()
  username := scanner.Text()

  conn := initConnection(username)

  rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(SCREEN_W, SCREEN_H, "Gooey")
  rl.SetTargetFPS(60)

  for !rl.WindowShouldClose() {
    rl.BeginDrawing()
    rl.ClearBackground(rl.Black)

    rl.EndDrawing()
  }

  // for {
  //   scanner.Scan()
  //   msg := scanner.Text()
  //   sendRegularMesage(conn, msg, myID)
  //
  //   if err != nil { panic(err) }
  // }

  rl.CloseWindow()
  conn.Close()
}
