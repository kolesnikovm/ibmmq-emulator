package mq

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
)

type tsh struct {
	StructID     string
	MQSegmLen    int
	SegmType     byte
	HeaderLength int
}

type tshm struct {
	StructID  []byte `offset:"0", length:"4"`
	MQSegmLen []byte `offset:"4", length:"4"`
	ConversID []byte `offset:"8", length:"4"`
	RequestID []byte `offset:"12", length:"4"`
	ByteOrder []byte `offset:"16", length:"1"`
	SegmType  []byte `offset:"17", length:"1"`
	CtlFlag1  []byte `offset:"18", length:"1"`
	CtlFlag2  []byte `offset:"19", length:"1"`
	LUWIdent  []byte `offset:"20", length:"8"`
	Encoding  []byte `offset:"28", length:"4"`
	CCSID     []byte `offset:"32", length:"2"`
	Reserved  []byte `offset:"34", length:"2"`
}

type apiHeader struct {
	ReplyLen   []byte `offset:"0", length:"4"`
	ComplCode  []byte `offset:"4", length:"4"`
	ReasonCode []byte `offset:"8", length:"4"`
	ObjectHdl  []byte `offset:"12", length:"4"`
}

type mqConn struct {
	QMgr        []byte `offset:"0", length:"48"`
	ApplName    []byte `offset:"48", length:"28"`
	ApplType    []byte `offset:"76", length:"4"`
	AccToken    []byte `offset:"80", length:"32"`
	Options     []byte `offset:"112", length:"4"`
	XOptions    []byte `offset:"116", length:"4"`
	FConnOption []byte `offset:"120", length:"212"`
}

type objectDescriptor struct {
	StructID       []byte `offset:"0" length:"4"`
	Version        []byte `offset:"4" length:"4"`
	ObjType        []byte `offset:"8" length:"4"`
	ObjName        []byte `offset:"12" length:"48"`
	ObjQMgr        []byte `offset:"60" length:"48"`
	DynQName       []byte `offset:"108" length:"48"`
	AltUserID      []byte `offset:"156" length:"12"`
	NbrRecord      []byte `offset:"168" length:"4"`
	KnownDestCount []byte `offset:"172" length:"4"`
	UnknownDestCnt []byte `offset:"176" length:"4"`
	InvalidDestCnt []byte `offset:"180" length:"4"`
	OffsetOR       []byte `offset:"184" length:"4"`
	OffsetRR       []byte `offset:"188" length:"4"`
	AddressOR      []byte `offset:"192" length:"4"`
	AddressRR      []byte `offset:"196" length:"4"`
}

const (
	INITIAL_DATA           = byte(0x01)
	USER_DATA              = byte(0x08)
	MQCONN, MQCONN_REPLY   = byte(0x81), byte(0x91)
	MQOPEN, MQOPEN_REPLY   = byte(0x83), byte(0x93)
	MQINQ, MQINQ_REPLY     = byte(0x89), byte(0x99)
	MQCLOSE, MQCLOSE_REPLY = byte(0x84), byte(0x94)
	SOCKET_ACTION          = byte(0x0c)
	SPI, SPI_REPLY         = byte(0x8c), byte(0x9c)
	MQPUT, MQPUT_REPLY     = byte(0x86), byte(0x96)
	REQUEST_MSGS           = byte(0x0e)
	ASYNC_MESSAGE          = byte(0x0d)
	NOTIFICATION           = byte(0x0f)
	MQCMIT, MQCMIT_REPLY   = byte(0x8a), byte(0x9a)
	MQDISC, MQDISC_REPLY   = byte(0x82), byte(0x92)
)

var (
	msgTypes = map[byte]byte{
		INITIAL_DATA: INITIAL_DATA,
		MQCONN:       MQCONN_REPLY,
		MQOPEN:       MQOPEN_REPLY,
		MQINQ:        MQINQ_REPLY,
		MQCLOSE:      MQCLOSE_REPLY,
		SPI:          SPI_REPLY,
		MQPUT:        MQPUT_REPLY,
		MQCMIT:       MQCMIT_REPLY,
		MQDISC:       MQDISC_REPLY,
	}

	userID, appType, appName, qMgr []byte //перенести в контекст
)

func HandleMessage(msg []byte) (response []byte) {
	tshType := msg[:4]
	fmt.Printf("tshType: %s\n", tshType)

	var msgType byte
	var tshmRs tshm

	switch string(tshType) {
	case "TSH ":
		msgType = msg[9]
	case "TSHM":
		tshmRs = tshm{
			StructID:  msg[0:4],
			ConversID: msg[8:12],
			RequestID: msg[12:16],
			ByteOrder: msg[16:17],
			SegmType:  []byte{msgTypes[msg[17]]},
			CtlFlag1:  msg[18:19],
			CtlFlag2:  msg[19:20],
			LUWIdent:  msg[20:28],
			Encoding:  msg[28:32],
			CCSID:     msg[32:34],
			Reserved:  msg[34:36],
		}
		msgType = msg[17]
		fmt.Printf("SegmType: 0x%x\n", tshmRs.SegmType)
	case "TSHC":
		msgType = msg[9]
	default:
		fmt.Printf("Unknown TSH type: %s\n", tshType)
	}

	switch msgType {
	case INITIAL_DATA:
		switch string(tshType) {
		case "TSH ":
			response, _ = hex.DecodeString("545348200000010c0201020000000000000000001101000033030000494420200d26009400000000ec7f000000004000000000004445562e41444d494e2e535652434f4e4e20202051003303514d312020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202c0100000000000000ff00ffffffffffffffffffffffffffffff0000000000000a000000000800009a030000150000003c0000004d514d4d3039303130323030514d315f323031392d31322d30355f31312e31372e30312020202020202020202020202020202020202020202020202001000000ffffffffffffffffffffffffffffffff051b5465bff9cf4adaba0a4c")
		case "TSHM":
			response, _ = hex.DecodeString("5453484d0000011400000001000000000201000000000000000000002202000033030000494420200d26000000000000ec7f000000004000000000004445562e41444d494e2e535652434f4e4e20202051003303514d312020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202c0100000000000000ff00ffffffffffffffffffffffffffffff0000000000000a000000000000009a030000150000003c0000004d514d4d3039303130323030514d315f323031392d31322d30355f31312e31372e30312020202020202020202020202020202020202020202020202001000000ffffffffffffffffffffffffffffffff051b5465bff9cf4adaba0a4c")
		}
	case USER_DATA:
		userID = msg[40:52]
	case MQCONN:
		tshmRs.MQSegmLen = []byte{0x00, 0x00, 0x01, 0x80}
		apiHeader := apiHeader{
			ReplyLen:   []byte{0x00, 0x00, 0x01, 0x78},
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  []byte{0x00, 0x00, 0x00, 0x00},
		}
		appType = msg[128:132]
		appName = msg[100:128]
		qMgr = msg[52:100]
		mqConn := mqConn{
			QMgr:        msg[52:100],
			ApplName:    msg[100:128],
			ApplType:    msg[128:132],
			AccToken:    msg[132:164],
			Options:     msg[164:168],
			XOptions:    msg[168:172],
			FConnOption: decodeString("46434e4f0200000001000000414d5143514d312020202020202020206445ea5d02b494240000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		}
		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
		response = append(response, getBytes(mqConn)...)
	case MQOPEN:
		tshmRs.MQSegmLen = []byte{0x00, 0x00, 0x01, 0x54}
		apiHeader := apiHeader{
			ReplyLen:   []byte{0x00, 0x00, 0x01, 0x4c},
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  []byte{0x02, 0x00, 0x00, 0x00},
		}
		objectDescriptor := objectDescriptor{
			StructID:       msg[52:56],
			Version:        msg[56:60],
			ObjType:        msg[60:64],
			ObjName:        msg[64:112],
			ObjQMgr:        msg[112:160],
			DynQName:       msg[160:208],
			AltUserID:      msg[208:220],
			NbrRecord:      msg[220:224],
			KnownDestCount: []byte{0x01, 0x00, 0x00, 0x00},
			UnknownDestCnt: msg[228:232],
			InvalidDestCnt: msg[232:236],
			OffsetOR:       msg[236:240],
			OffsetRR:       msg[240:244],
			AddressOR:      msg[244:248],
			AddressRR:      msg[248:252],
		}
		mqOpen := mqOpen{
			Options: []byte{0x20, 0x20, 0x00, 0x00},
		}
		hiddenField, _ := hex.DecodeString("2020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020205400000000000000")
		fopa := fopa{
			StructID:        msg[256:260],
			Version:         msg[260:264],
			Length:          msg[264:268],
			DefPersistence:  msg[268:272],
			DefPutRespType:  []byte{0xff, 0xff, 0xff, 0xff},
			DefReadAhead:    []byte{0xff, 0xff, 0xff, 0xff},
			PropertyControl: []byte{0xff, 0xff, 0xff, 0xff},
			hiddenField:     hiddenField,
		}
		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
		response = append(response, getBytes(objectDescriptor)...)
		response = append(response, getBytes(mqOpen)...)
		response = append(response, getBytes(fopa)...)
	case MQINQ:
		mqInc := handleMqInc(msg[52:])
		msgLength := 36 + 16 + len(mqInc)
		msgLengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLengthBytes, uint32(msgLength))

		replyLen := msgLength - 8
		replyLenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(replyLenBytes, uint32(replyLen))

		tshmRs.MQSegmLen = msgLengthBytes
		apiHeader := apiHeader{
			ReplyLen:   replyLenBytes,
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  []byte{0x02, 0x00, 0x00, 0x00},
		}
		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
		response = append(response, mqInc...)
	case MQCLOSE:
		tshmRs.MQSegmLen = []byte{0x00, 0x00, 0x00, 0x34}
		apiHeader := apiHeader{
			ReplyLen:   []byte{0x00, 0x00, 0x00, 0x2c},
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  []byte{0xff, 0xff, 0xff, 0xff},
		}
		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
	case SOCKET_ACTION:
		if bytes.Compare(msg[36:40], []byte{0x02, 0x00, 0x00, 0x00}) == 0 {
			return nil
		}
		tshmRs := tshm{
			StructID:  []byte{0x54, 0x53, 0x48, 0x4d},
			MQSegmLen: []byte{0x00, 0x00, 0x00, 0x38},
			ConversID: []byte{0x00, 0x00, 0x00, 0x03},
			RequestID: msg[32:36],
			ByteOrder: msg[8:9],
			SegmType:  []byte{SOCKET_ACTION},
			CtlFlag1:  []byte{0x00},
			CtlFlag2:  []byte{0x00},
			LUWIdent:  []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			Encoding:  []byte{0x22, 0x02, 0x00, 0x00},
			CCSID:     []byte{0x33, 0x03},
			Reserved:  []byte{0x00, 0x00},
		}
		socketAction := socketAction{
			ConversID: msg[28:32],
			RequestID: msg[32:36],
			Type:      []byte{0x08, 0x00, 0x00, 0x00},
			Param1:    []byte{0x15, 0x00, 0x00, 0x00},
			Param2:    []byte{0x3d, 0x00, 0x00, 0x00},
		}

		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(socketAction)...)
	case SPI:
		spi := handleSpi(msg[52:])
		msgLength := 36 + 16 + len(spi)
		msgLengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLengthBytes, uint32(msgLength))

		replyLen := msgLength - 8
		replyLenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(replyLenBytes, uint32(replyLen))

		tshmRs.MQSegmLen = msgLengthBytes

		objectHdl := []byte{0x00, 0x00, 0x00, 0x00}
		if bytes.Compare(msg[52:56], []byte{0x0c, 0x00, 0x00, 0x00}) == 0 {
			lpiVersion := msg[96:100]
			switch {
			case bytes.Compare(lpiVersion, []byte{0x10, 0x20, 0x00, 0x00}) == 0:
				objectHdl = []byte{0x02, 0x00, 0x00, 0x00}
			case bytes.Compare(lpiVersion, []byte{0x29, 0x20, 0x00, 0x00}) == 0:
				objectHdl = []byte{0x04, 0x00, 0x00, 0x00}
			}
		}

		apiHeader := apiHeader{
			ReplyLen:   replyLenBytes,
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  objectHdl, //на каждое открытие свой номер
		}

		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
		response = append(response, spi...)
	case MQPUT:
		mqPut := handleMqPut(msg[52:], userID, appType, appName, qMgr)
		msgLength := 36 + 16 + len(mqPut)
		msgLengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLengthBytes, uint32(msgLength))

		replyLen := msgLength - 8
		replyLenBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(replyLenBytes, uint32(replyLen))

		tshmRs.MQSegmLen = msgLengthBytes

		apiHeader := apiHeader{
			ReplyLen:   replyLenBytes,
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  msg[48:52],
		}

		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
		response = append(response, mqPut...)
	case REQUEST_MSGS:
		asyncMsg := handleRequestMsg(msg[36:], userID, appType, appName, qMgr)

		msgLength := 36 + len(asyncMsg)
		msgLengthBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(msgLengthBytes, uint32(msgLength))

		tshmRs.MQSegmLen = msgLengthBytes
		tshmRs.RequestID = []byte{0x00, 0x00, 0x00, 0x01}
		tshmRs.SegmType = []byte{ASYNC_MESSAGE}
		tshmRs.CtlFlag1 = []byte{0x30}
		tshmRs.Encoding = []byte{0x22, 0x02, 0x00, 0x01}

		response = append(response, getBytes(tshmRs)...)
		response = append(response, asyncMsg...)
	case MQCMIT:
		tshmRs.MQSegmLen = []byte{0x00, 0x00, 0x00, 0x34}
		apiHeader := apiHeader{
			ReplyLen:   []byte{0x00, 0x00, 0x00, 0x2c},
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  []byte{0x00, 0x00, 0x00, 0x00},
		}

		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
	case MQDISC:
		tshmRs.MQSegmLen = []byte{0x00, 0x00, 0x00, 0x34}
		apiHeader := apiHeader{
			ReplyLen:   []byte{0x00, 0x00, 0x00, 0x2c},
			ComplCode:  []byte{0x00, 0x00, 0x00, 0x00},
			ReasonCode: []byte{0x00, 0x00, 0x00, 0x00},
			ObjectHdl:  []byte{0x00, 0x00, 0x00, 0x00},
		}

		response = append(response, getBytes(tshmRs)...)
		response = append(response, getBytes(apiHeader)...)
	}

	return response
}

func decodeString(stg string) []byte {
	bytes, _ := hex.DecodeString(stg)
	return bytes
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
