package main

import (
	"github.com/gansidui/gotcp"
	"log"
	"net"
	"time"

	"demo/gotcp_test/echo"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Runmain() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:9905")
	checkError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer func() {
		err = conn.Close()
		if err != nil {

		}
	}()

	protocol := &echo.EchoMsgProtocol{}

	for i := 0; i < 3; i++ {
		var n int
		n, err = conn.Write(echo.NewEchoMsgPacket([]byte("hello"), false).Serialize())
		if err != nil {

		}
		var p gotcp.Packet
		p, err = protocol.ReadPacket(conn)
		if err == nil {
			npack := p.(*echo.EchoMsgPacket)
			log.Printf("Server reply:[%v] [%v]\n", npack.GetLength(), string(npack.GetBody()))
		}

		time.Sleep(time.Second * time.Duration(1))
	}
}

func main() {
	Runmain()
}
