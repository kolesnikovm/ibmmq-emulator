package mq

import (
	"bytes"
	"log"
	"encoding/binary"
)

func handleMqClose(msg []byte) (response []byte) {
	log.Printf("[INFO] M: MQCLOSE, C: %d, R: %d, Hdl: %d\n", binary.BigEndian.Uint32(msg[8:12]), binary.BigEndian.Uint32(msg[12:16]), binary.LittleEndian.Uint32(msg[48:52]))

	apiHeader := apiHeader{
		ReplyLen:   []byte{0x00, 0x00, 0x00, 0x2c},
		ComplCode:  MQCC_OK,
		ReasonCode: MQRC_NONE,
		ObjectHdl:  getObjectHdl(msg[48:52]),
	}
	response = append(response, getBytes(apiHeader)...)

	return response
}

func getObjectHdl(requestHdl []byte) (responseHdl []byte) {
	switch {
	case bytes.Compare(requestHdl, []byte{0x02, 0x00, 0x00, 0x00}) == 0:
		responseHdl = []byte{0xff, 0xff, 0xff, 0xff}
	case bytes.Compare(requestHdl, []byte{0x04, 0x00, 0x00, 0x00}) == 0:
		responseHdl = []byte{0x04, 0x00, 0x00, 0x00}
	default:
		log.Printf("[WARN] unknown ObjectHdl: %x\n", requestHdl)
		responseHdl = []byte{0xff, 0xff, 0xff, 0xff}
	}

	return responseHdl
}
