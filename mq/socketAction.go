package mq

import (
	"encoding/binary"
	"log"
)

type socketAction struct {
	ConversID []byte
	RequestID []byte
	Type      []byte
	Param1    []byte
	Param2    []byte
}

func handleSocketAction(msg []byte) (response []byte) {
	log.Printf("[INFO] M: SOCKET_ACTION, C: %d, R: %d, P1: %d, P2: %d\n", binary.LittleEndian.Uint32(msg[28:32]), binary.LittleEndian.Uint32(msg[32:36]), binary.LittleEndian.Uint32(msg[40:44]), binary.LittleEndian.Uint32(msg[44:48]))

	tshm := tshm{
		StructID:  []byte(TSHM),
		MQSegmLen: []byte{0x00, 0x00, 0x00, 0x38},
		ConversID: []byte{msg[31], msg[30], msg[29], msg[28]},
		RequestID: []byte{msg[35], msg[34], msg[33], msg[32]},
		ByteOrder: msg[8:9],
		SegmType:  []byte{SOCKET_ACTION},
		CtlFlag1:  []byte{0x00},
		CtlFlag2:  []byte{0x00},
		LUWIdent:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Encoding:  []byte{0x22, 0x02, 0x00, 0x00},
		CCSID:     []byte{0x33, 0x03},
		Reserved:  []byte{0x00, 0x00},
	}
	response = append(response, getBytes(tshm)...)

	socketAction := socketAction{
		ConversID: msg[28:32],
		RequestID: msg[32:36],
		Type:      []byte{0x08, 0x00, 0x00, 0x00},
		Param1:    []byte{0x15, 0x00, 0x00, 0x00}, // глобальные счетчики
		Param2:    []byte{0x3d, 0x00, 0x00, 0x00},
	}

	response = append(response, getBytes(socketAction)...)

	return response
}
