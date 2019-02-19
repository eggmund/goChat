package cliTools

import "net"

type CliData struct {
  Conn *net.TCPConn
  ID *CliID
}
