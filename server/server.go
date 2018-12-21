package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gansidui/gotcp"

	"gotcp_test/echo"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Callback struct {
}

func (this *Callback) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr()
	c.PutExtraData(addr)

	fmt.Println("OnConnect:", addr)

	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*echo.MsgPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", packet.GetLength(), string(packet.GetBody()))
	c.AsyncWritePacket(echo.NewMsgPacket(p.Serialize(), true), time.Second)
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	fmt.Println("OnClose", c.GetExtraData())
}

func Runmain() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":9905")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	config := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(config, &Callback{}, &echo.MsgProtocol{})

	go srv.Start(listener, time.Second*time.Duration(5))

	fmt.Println("listening:", listener.Addr())

	chSig := make(chan os.Signal)

	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Signal:", <-chSig)

	srv.Stop()
}

func main() {
	Runmain()
}