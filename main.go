package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"./mq"
)

type tshm struct {
	StructID  []byte
	MQSegmLen []byte
	ConversID []byte
	RequestID []byte
	ByteOrder []byte
	SegmType  []byte
	CtlFlag1  []byte
	CtlFlag2  []byte
	LUWIdent  []byte
	Encoding  []byte
	CCSID     []byte
	Reserved  []byte
}

func getBytes(msgPart interface{}) (bytes []byte) {
	v := reflect.ValueOf(msgPart)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsNil() {
			continue
		}

		bytes = append(bytes, v.Field(i).Bytes()...)
	}

	return bytes
}

func handleConn(conn net.Conn) {
	fmt.Println("Handling new connection...")

	defer func() {
		fmt.Println("Closing connection...")
		conn.Close()
	}()

	//create context
	// ctx := make(map[string][]byte)

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))

		messageStart := make([]byte, 8)
		_, err := io.ReadFull(conn, messageStart)
		if err != nil {
			return
		}

		length := int(binary.BigEndian.Uint32(messageStart[4:]))

		messageEnd := make([]byte, length-len(messageStart))
		_, err = io.ReadFull(conn, messageEnd)
		if err != nil {
			return
		}

		message := append(messageStart, messageEnd...)

		response := mq.HandleMessage(message)
		conn.Write(response)

		if len(response) > 17 && response[17] == mq.ASYNC_MESSAGE {
			var notification []byte
			tshm := tshm{
				StructID:  response[0:4],
				MQSegmLen: []byte{0x00, 0x00, 0x00, 0x34},
				ConversID: response[8:12],
				RequestID: []byte{0x00, 0x00, 0x00, 0x03},
				ByteOrder: response[16:17],
				SegmType:  []byte{mq.NOTIFICATION},
				CtlFlag1:  []byte{0x10},
				CtlFlag2:  response[19:20],
				LUWIdent:  response[20:28],
				Encoding:  response[28:32],
				CCSID:     response[32:34],
				Reserved:  response[34:36],
			}
			notif, _ := hex.DecodeString("01000000040000000b00000000000000")
			notification = append(notification, getBytes(tshm)...)
			notification = append(notification, notif...)

			conn.Write(notification)
		}
	}
}

func getMessageLength(conn net.Conn) (int, []byte, error) {
	message := make([]byte, 8)
	_, err := io.ReadFull(conn, message)
	if err != nil {
		return 0, nil, err
	}
	length := int(binary.BigEndian.Uint32(message[4:]))

	return length, message, nil
}

func main() {
	fmt.Println("Launching server...")

	ln, _ := net.Listen("tcp", ":1414")

	for {
		conn, _ := ln.Accept()
		go handleConn(conn)
	}
}
