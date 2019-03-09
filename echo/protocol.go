package echo

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/gansidui/gotcp"
)

type EchoMsgPacket struct {
	buff []byte
}

type EchoMsgProtocol struct {
}

func (emp *EchoMsgPacket) Serialize() []byte {
	return emp.buff
}

func (emp *EchoMsgPacket) GetLength() uint32 {
	return binary.BigEndian.Uint32(emp.buff[0:4])
}

func (emp *EchoMsgPacket) GetBody() []byte {
	return emp.buff[4:]
}

func NewEchoMsgPacket(buff []byte, hasLengthField bool) *EchoMsgPacket {
	p := &EchoMsgPacket{}
	if hasLengthField {
		p.buff = buff
	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff, uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

func (emp *EchoMsgProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	lengthBytes := make([]byte, 4)
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}

	var (
		length uint32
	)

	if length = binary.BigEndian.Uint32(lengthBytes); length > 1024 {
		return nil, fmt.Errorf("the size of packet is larger than the limit")
	}

	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return nil, err
	}

	return NewEchoMsgPacket(buff, true), nil
}
