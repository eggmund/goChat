package msg

import (
  "cliTools"
)

/* Type: 0 = Regular message
         1 = Handshake
         2 = Ping
*/
type Message struct {
  Type uint8
  Content string
  Author cliTools.CliID
}

func NewMessage(t uint8, content string, author *cliTools.CliID) Message {
  return Message{
    Type: t,
    Content: content,
    Author: *author,
  }
}
