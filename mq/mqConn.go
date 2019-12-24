package mq

import (
	"encoding/binary"
	"log"
	"strings"
)

type mqConn struct {
	QMgr        []byte `offset:"0", length:"48"`
	ApplName    []byte `offset:"48", length:"28"`
	ApplType    []byte `offset:"76", length:"4"`
	AccToken    []byte `offset:"80", length:"32"`
	Options     []byte `offset:"112", length:"4"`
	XOptions    []byte `offset:"116", length:"4"`
	FConnOption []byte `offset:"120", length:"212"`
}

func handleMqConn(msg []byte) (response []byte) {
	log.Printf("[INFO] M: MQCONN, C: %d, R: %d, A: %s, Q: %s\n", binary.BigEndian.Uint32(msg[8:12]), binary.BigEndian.Uint32(msg[12:16]), strings.TrimSpace(string(msg[100:128])), strings.TrimSpace(string(msg[52:100])))

	mqConn := mqConn{
		QMgr:        msg[52:100],
		ApplName:    msg[100:128],
		ApplType:    msg[128:132],
		AccToken:    msg[132:164],
		Options:     msg[164:168],
		XOptions:    msg[168:172],
		FConnOption: decodeString("46434e4f0200000001000000414d5143514d312020202020202020206445ea5d02b494240000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
	}
	response = append(response, getBytes(mqConn)...)

	return response
}
