package mq

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/rand"
)

type asyncMessage struct {
	Version   []byte
	Handle    []byte
	MsgIndex  []byte
	GlbMsgIdx []byte
	SegLength []byte
	SegIndex  []byte
	SelIndex  []byte
	ReasonCod []byte
	TotMsgLen []byte
	ActMsgLen []byte
	MsgToken  []byte
	Status    []byte
	ResolQNLn []byte
	ResolQNme []byte
	Padding   []byte
}

type asyncMessageDescriptor struct {
	StructID   []byte
	Version    []byte
	Report     []byte
	MsgType    []byte
	Expiry     []byte
	Feedback   []byte
	Encoding   []byte
	CCSID      []byte
	Format     []byte
	Priority   []byte
	Persist    []byte
	MsgID      []byte
	CorrelID   []byte
	BackoCnt   []byte
	ReplyToQ   []byte
	ReplToQMgr []byte
	UserID     []byte
	AccntTok   []byte
	AppIDDAta  []byte
	PutAppTyp  []byte
	PutAppName []byte
	PutDatGMT  []byte
	PutTimGMT  []byte
	AppOriDat  []byte
	GroupID    []byte
	MsgSeqNum  []byte
	Offset     []byte
	MsgFlags   []byte
	OrigLen    []byte
}

type rulesFormattingHeader struct {
	StructID    []byte
	Version     []byte
	Length      []byte
	Encoding    []byte
	CCSID       []byte
	Format      []byte
	Flags       []byte
	NmeValCCSID []byte
	McdValue    []byte
	JmsValue    []byte
	UsrValue    []byte
}

type data struct {
	Value []byte
}

func handleRequestMsg(msg, userID, appType, appName, qMgr []byte) (response []byte) {
	log.Printf("[INFO] M: REQUEST_MESSAGE, C: %d, R: %d\n", binary.BigEndian.Uint32(msg[8:12]), binary.BigEndian.Uint32(msg[12:16]))

	msgToken, _ := hex.DecodeString("6345ea5d410000005f00000000000000")

	queueName := []byte("DEV.QUEUE.1 ") //нужен ли в конце пробел
	queueNameLen := byte(len(queueName))

	asyncMessage := asyncMessage{
		Version:   msg[36:40],
		Handle:    msg[40:44],
		MsgIndex:  []byte{0x01, 0x00, 0x00, 0x00},
		GlbMsgIdx: []byte{0x01, 0x00, 0x00, 0x00},
		SegLength: make([]byte, 4),
		SegIndex:  []byte{0x00, 0x00},
		SelIndex:  msg[76:78],
		ReasonCod: MQRC_NONE,
		TotMsgLen: make([]byte, 4),
		ActMsgLen: make([]byte, 4),
		MsgToken:  msgToken,
		Status:    []byte{0x01, 0x00},
		ResolQNLn: []byte{queueNameLen},
		ResolQNme: queueName,
		Padding:   []byte{0x00},
	}

	msgID := make([]byte, 24)
	rand.Read(msgID)

	correlID := make([]byte, 24)

	replyQ := make([]byte, 48)
	for i := 0; i < 48; i++ {
		replyQ[i] = byte(0x20)
	}

	accntToken := make([]byte, 32)

	appIDDAta := make([]byte, 32)
	for i := 0; i < 32; i++ {
		appIDDAta[i] = byte(0x20)
	}

	groupID := make([]byte, 24)

	asyncMessageDescriptor := asyncMessageDescriptor{
		StructID:   []byte{0x4d, 0x44, 0x20, 0x20},
		Version:    []byte{0x02, 0x00, 0x00, 0x00},
		Report:     []byte{0x00, 0x00, 0x00, 0x00},
		MsgType:    []byte{0x08, 0x00, 0x00, 0x00},
		Expiry:     []byte{0xff, 0xff, 0xff, 0xff},
		Feedback:   []byte{0x00, 0x00, 0x00, 0x00},
		Encoding:   []byte{0x11, 0x01, 0x00, 0x00},
		CCSID:      []byte{0xb8, 0x04, 0x00, 0x00},
		Format:     []byte{0x4d, 0x51, 0x48, 0x52, 0x46, 0x32, 0x20, 0x20},
		Priority:   []byte{0x04, 0x00, 0x00, 0x00},
		Persist:    []byte{0x01, 0x00, 0x00, 0x00},
		MsgID:      msgID,
		CorrelID:   correlID,
		BackoCnt:   []byte{0x00, 0x00, 0x00, 0x00},
		ReplyToQ:   replyQ,
		ReplToQMgr: qMgr,
		UserID:     userID,
		AccntTok:   accntToken,
		AppIDDAta:  appIDDAta,
		PutAppTyp:  appType,
		PutAppName: appName,
		PutDatGMT:  []byte{0x32, 0x30, 0x31, 0x39, 0x31, 0x32, 0x30, 0x36},
		PutTimGMT:  []byte{0x31, 0x32, 0x34, 0x34, 0x31, 0x33, 0x35, 0x37},
		AppOriDat:  []byte{0x20, 0x20, 0x20, 0x20},
		GroupID:    groupID,
		MsgSeqNum:  []byte{0x01, 0x00, 0x00, 0x00},
		Offset:     []byte{0x00, 0x00, 0x00, 0x00},
		MsgFlags:   []byte{0x00, 0x00, 0x00, 0x00},
		OrigLen:    []byte{0xff, 0xff, 0xff, 0xff},
	}

	mcdValue, _ := hex.DecodeString("000000203c6d63643e3c4d73643e6a6d735f746578743c2f4d73643e3c2f6d63643e2020")
	jmcValue, _ := hex.DecodeString("000000503c6a6d733e3c4473743e71756575653a2f2f2f4445562e51554555452e313c2f4473743e3c546d733e313537353633363235333537343c2f546d733e3c446c763e323c2f446c763e3c2f6a6d733e2020")
	usrValue, _ := hex.DecodeString("000000703c7573723e3c6573666c5f6d6574686f644e616d653e636c69656e744372656174653c2f6573666c5f6d6574686f644e616d653e3c7372635f73797374656d49443e4452544c3c2f7372635f73797374656d49443e3c76657273696f6e3e323c2f76657273696f6e3e3c2f7573723e20")

	rulesFormattingHeader := rulesFormattingHeader{
		StructID:    []byte{0x52, 0x46, 0x48, 0x20},
		Version:     []byte{0x00, 0x00, 0x00, 0x02},
		Length:      make([]byte, 4),
		Encoding:    []byte{0x00, 0x00, 0x01, 0x11},
		CCSID:       []byte{0x00, 0x00, 0x04, 0xb8},
		Format:      []byte{0x4d, 0x51, 0x53, 0x54, 0x52, 0x20, 0x20, 0x20},
		Flags:       []byte{0x00, 0x00, 0x00, 0x00},
		NmeValCCSID: []byte{0x00, 0x00, 0x04, 0xb8},
		McdValue:    mcdValue,
		JmsValue:    jmcValue,
		UsrValue:    usrValue,
	}
	rfhLength := make([]byte, 4)
	binary.BigEndian.PutUint32(rfhLength, uint32(len(getBytes(rulesFormattingHeader))))
	rulesFormattingHeader.Length = rfhLength

	payload := "lox bus alesha"
	normPayload := append([]byte(payload), make([]byte, 4-len(payload)%4)...)

	payloadLength := make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadLength, uint32(len(getBytes(rulesFormattingHeader))+len(payload)))
	asyncMessage.SegLength = payloadLength
	asyncMessage.TotMsgLen = payloadLength
	asyncMessage.ActMsgLen = payloadLength

	response = append(response, getBytes(asyncMessage)...)
	response = append(response, getBytes(asyncMessageDescriptor)...)
	response = append(response, getBytes(rulesFormattingHeader)...)

	data := data{
		Value: normPayload,
	}
	response = append(response, getBytes(data)...)

	return response
}
