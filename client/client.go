package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"gotcp_test/echo"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Runmain() {
	tcpAddr, error := net.ResolveTCPAddr("tcp4", "127.0.0.1:9905")
	checkError(error)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	defer conn.Close()

	protocol := &echo.EchoMsgProtocol{}

	for i := 0; i < 3; i++ {
		conn.Write(echo.NewEchoMsgPacket([]byte("hello"), false).Serialize())

		p, err := protocol.ReadPacket(conn)
		if err == nil {
			npack := p.(*echo.EchoMsgPacket)
			fmt.Printf("Server reply:[%v] [%v]\n", npack.GetLength(), string(npack.GetBody()))
		}

		time.Sleep(time.Second * time.Duration(2))
	}
}

func main() {
	Runmain()
}
