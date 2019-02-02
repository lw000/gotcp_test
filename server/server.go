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

	"demo/gotcp_test/echo"
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

	log.Println("OnConnect:", addr)

	return true
}

func (this *Callback) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*echo.EchoMsgPacket)
	fmt.Printf("OnMessage:[%v] [%v]\n", packet.GetLength(), string(packet.GetBody()))
	err := c.AsyncWritePacket(echo.NewEchoMsgPacket(p.Serialize(), true), time.Second)
	if err != nil {

	}
	return true
}

func (this *Callback) OnClose(c *gotcp.Conn) {
	log.Println("OnClose", c.GetExtraData())
}

func Runmain() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":9905")
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	conf := &gotcp.Config{
		PacketSendChanLimit:    20,
		PacketReceiveChanLimit: 20,
	}
	srv := gotcp.NewServer(conf, &Callback{}, &echo.EchoMsgProtocol{})

	go srv.Start(listener, time.Second*time.Duration(1))

	log.Println("listening:", listener.Addr())

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	log.Println("signal:", <-c)
	close(c)
	srv.Stop()
}

func Runwait() {

}

func main() {
	Runmain()
}
