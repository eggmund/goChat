package main

import (
  "fmt"
  "net"
  "os"
  "bufio"
  "encoding/json"
  "github.com/gotk3/gotk3/gtk"

  "msg"
  "cliTools"
)

var (
  SCREEN_W int32 = 1000
  SCREEN_H int32 = 600
  messages []*msg.Message
  msgTextView *gtk.TextView

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
    check(err)
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
        messages = append(messages, &m)
        updateTextBuffer(msgTextView, messages)
      }
    }
  }
}

func convDataToMsg(data []byte) msg.Message {
  var m msg.Message
  err := json.Unmarshal(data, &m)
  check(err)

  return m
}

func sendRegularMesage(c *net.TCPConn, content string, author *cliTools.CliID) {
  byt, err := json.Marshal( msg.NewMessage(0, content, author) )
  check(err)

  (*c).Write(byt)
}

func getMyID(c *net.TCPConn, username string) *cliTools.CliID {
  byt, err := json.Marshal( msg.NewMessage(1, username, nil) )
  check(err)
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
  check(err)

  conn, err := net.DialTCP("tcp", nil, tcpAddr)
  check(err)
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

  gtk.Init(nil)

  win := setupWindow("goChat")
  grid, err := gtk.GridNew()
  check(err)

  grid.SetOrientation(gtk.ORIENTATION_VERTICAL)


  msgScrl, err := gtk.ScrolledWindowNew(nil, nil)
  check(err)

  msgTextView, err = gtk.TextViewNew()
  check(err)
  msgTextView.SetEditable(false)

  msgScrl.Add(msgTextView)

  winW, winH := win.GetSize()
  msgScrl.SetSizeRequest(winW, winH-30)

  buffer := getTvBuffer(msgTextView)
  buffer.SetText(getTextBufferFromMessages(messages))

  msgEntry, err := gtk.EntryNew()
  check(err)
  msgEntry.Connect("activate", func() {
    buff, err := msgEntry.GetBuffer()
    check(err)

    str, err := buff.GetText()
    check(err)
    sendRegularMesage(conn, str, myID)
    buff.SetText("")
  })

  grid.Add(msgScrl)
  grid.Add(msgEntry)

  win.Add(grid)
  win.ShowAll()
  gtk.Main()

  conn.Close()
}
