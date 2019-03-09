// gotcp_test project main.go
package main

import (
	"bytes"
	"fmt"

	// "gotcp_test/client"
	// "gotcp_test/server"
	// "os"

	"encoding/binary"
)

func Int64ToBytes(v int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, v)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buffer.Bytes()
}

func Int32ToBytes(v uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}

func main() {
	// args := os.Args
	// if args[0] == "1" {
	// 	server.Runmain()
	// } else if args[0] == "0" {
	// 	client.Runmain()
	// } else {
	// 	fmt.Println("error")
	// }

	{
		buf := Int64ToBytes(1234)
		value := binary.BigEndian.Uint64(buf)
		fmt.Println(value)
	}
	{
		buf := Int32ToBytes(1234)
		value := binary.BigEndian.Uint32(buf)
		fmt.Println(value)
	}
}
