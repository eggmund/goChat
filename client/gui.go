package main

import (
  "fmt"
  "msg"
  "github.com/gotk3/gotk3/gtk"
)

func setupWindow(title string) *gtk.Window {
  win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
  check(err)

  win.SetTitle("goChat")
  win.Connect("destroy", func() {
    gtk.MainQuit()
  })

  win.SetDefaultSize(1000, 600)
  return win
}

func getTvBuffer(tv *gtk.TextView) *gtk.TextBuffer {
  buffer, err := tv.GetBuffer()
  check(err)
  return buffer
}

func getTextBufferFromMessages(msgs []*msg.Message) string {
  var out string
  for i := range msgs {
    if ms, ok := msgs[i].Content.(string); ok {
      out += fmt.Sprintf("%s:: %s\n", msgs[i].Author.Username, ms)+"\n"
    }
  }
  return out
}

func updateTextBuffer(tv *gtk.TextView, msgs []*msg.Message) {
  if len(msgs) > 100 {
    msgs = append(msgs[:0], msgs[1:]...)
  }

  buff := getTvBuffer(tv)
  buff.SetText(getTextBufferFromMessages(msgs))
  tv.QueueDraw()
}
