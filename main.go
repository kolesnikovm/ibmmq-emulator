package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"time"

	"./mq"
)

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
	log.Println("[INFO] ===== Handling new connection =====")

	defer func() {
		log.Println("[INFO] ===== Closing connection =====")
		conn.Close()
	}()

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
