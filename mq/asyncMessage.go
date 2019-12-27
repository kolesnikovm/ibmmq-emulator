package mq

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"math/rand"
	"strings"
	"time"
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

type nameValue struct {
	length []byte
	value  []byte
}

type data struct {
	Value []byte
}

func handleRequestMsg(msg []byte) (response []byte) {
	cid := binary.BigEndian.Uint32(msg[8:12])
	rid := binary.BigEndian.Uint32(msg[12:16])
	hdl := binary.LittleEndian.Uint32(msg[40:44])

	log.Printf("[INFO] M: REQUEST_MESSAGE, C: %d, R: %d, H: %d\n", cid, rid, hdl)

	msgToken, _ := hex.DecodeString("6345ea5d410000005f00000000000000")

	q := ctx.sessions[cid].hdls[hdl].queue
	qN := strings.TrimSpace(string(q.name))
	queueName := make([]byte, 0, len(qN)+len(qN)%2)
	queueName = append(queueName, []byte(qN)...)
	for i := len(queueName); i < cap(queueName); i++ {
		queueName = append(queueName, []byte(" ")...)
	}
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

	t := time.Now()

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
		ReplToQMgr: ctx.sessions[cid].qMgr,
		UserID:     ctx.userID,
		AccntTok:   accntToken,
		AppIDDAta:  appIDDAta,
		PutAppTyp:  ctx.sessions[cid].appType,
		PutAppName: ctx.sessions[cid].appName,
		PutDatGMT:  []byte(t.Format("20060102")),
		PutTimGMT:  []byte(t.Format("15040500")),
		AppOriDat:  []byte{0x20, 0x20, 0x20, 0x20},
		GroupID:    groupID,
		MsgSeqNum:  []byte{0x01, 0x00, 0x00, 0x00},
		Offset:     []byte{0x00, 0x00, 0x00, 0x00},
		MsgFlags:   []byte{0x00, 0x00, 0x00, 0x00},
		OrigLen:    []byte{0xff, 0xff, 0xff, 0xff},
	}

	mcdValue, _ := hex.DecodeString("000000203c6d63643e3c4d73643e6a6d735f746578743c2f4d73643e3c2f6d63643e2020")

	message := q.get()

	jmsValue := nameValue{
		length: getByteLength(len(message.jmsValue)),
		value:  message.jmsValue,
	}

	usrValue := nameValue{
		length: getByteLength(len(message.usrValue)),
		value:  message.usrValue,
	}

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
		JmsValue:    getBytes(jmsValue),
		UsrValue:    getBytes(usrValue),
	}
	rfhLength := make([]byte, 4)
	binary.BigEndian.PutUint32(rfhLength, uint32(len(getBytes(rulesFormattingHeader))))
	rulesFormattingHeader.Length = rfhLength

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
