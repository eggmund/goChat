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
  Content interface{}
  Author cliTools.CliID
}

func NewMessage(t uint8, content interface{}, author *cliTools.CliID) Message {
  if author == nil {
    author = &cliTools.CliID {
      IDnum: -1,
      Username: "none",
    }
  }
  return Message{
    Type: t,
    Content: content,
    Author: *author,
  }
}
