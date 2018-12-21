package echo

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/gansidui/gotcp"
)

type MsgPacket struct {
	buff []byte
}

type MsgProtocol struct {
}

func (this *MsgPacket) Serialize() []byte {
	return this.buff
}

func (this *MsgPacket) GetLength() uint32 {
	return binary.BigEndian.Uint32(this.buff[0:4])
}

func (this *MsgPacket) GetBody() []byte {
	return this.buff[4:]
}

func NewMsgPacket(buff []byte, hasLengthField bool) *MsgPacket {
	p := &MsgPacket{}
	if hasLengthField {
		p.buff = buff
	} else {
		p.buff = make([]byte, 4+len(buff))
		binary.BigEndian.PutUint32(p.buff, uint32(len(buff)))
		copy(p.buff[4:], buff)
	}

	return p
}

func (this *MsgProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)

	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return nil, err
	}

	if length = binary.BigEndian.Uint32(lengthBytes); length > 1024 {
		return nil, fmt.Errorf("the size of packet is larger than the limit")
	}

	buff := make([]byte, 4+length)
	copy(buff[0:4], lengthBytes)

	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return nil, err
	}

	return NewMsgPacket(buff, true), nil
}
