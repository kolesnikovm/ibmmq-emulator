package mq

import (
	"encoding/binary"
	"encoding/hex"
	"log"
	"regexp"
	"time"
)

type messageDescriptor struct {
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
}

type putMessageOptions struct {
	StructID   []byte
	Version    []byte
	Options    []byte
	Timeout    []byte
	Context    []byte
	KnDstCnt   []byte
	UkDstCnt   []byte
	InDstCnt   []byte
	ResQName   []byte
	ResQMgr    []byte
	NumRecs    []byte
	PMRFlag    []byte
	OffsetPMR  []byte
	OffsetRR   []byte
	AddressPMR []byte
	AddressRR  []byte
}

type mqput struct {
	DataLength []byte
}

func handleMqPut(msg []byte) (response []byte) {
	cid := binary.BigEndian.Uint32(msg[8:12])
	rid := binary.BigEndian.Uint32(msg[12:16])
	hdl := binary.LittleEndian.Uint32(msg[48:52])
	queue := getQueueName(msg)

	log.Printf("[INFO] M: MQPUT, C: %d, R: %d, H: %d, Q: %s\n", cid, hdl, rid, queue)

	t := time.Now()

	msgID, _ := hex.DecodeString("414d5120514d312020202020202020206445ea5d04b59424") //сделать вручную
	messageDescriptor := messageDescriptor{
		StructID:   msg[52:56],
		Version:    msg[56:60],
		Report:     msg[60:64],
		MsgType:    msg[64:68],
		Expiry:     msg[68:72],
		Feedback:   msg[72:76],
		Encoding:   msg[76:80],
		CCSID:      msg[80:84],
		Format:     msg[84:92],
		Priority:   msg[92:96],
		Persist:    msg[96:100],
		MsgID:      msgID,
		CorrelID:   msg[124:148],
		BackoCnt:   msg[148:152],
		ReplyToQ:   msg[152:200],
		ReplToQMgr: msg[200:248],
		UserID:     ctx.userID,
		AccntTok:   msg[260:292],
		AppIDDAta:  msg[292:324],
		PutAppTyp:  ctx.sessions[cid].appType,
		PutAppName: ctx.sessions[cid].appName,
		PutDatGMT:  []byte(t.Format("20060102")),
		PutTimGMT:  []byte(t.Format("15040500")),
		AppOriDat:  []byte{0x20, 0x20, 0x20, 0x20},
	}
	response = append(response, getBytes(messageDescriptor)...)

	resQName := make([]byte, 0, 48)
	resQName = append(resQName, []byte(queue)...)
	for i := len(resQName); i < cap(resQName); i++ {
		resQName = append(resQName, []byte(" ")...)
	}

	putMessageOptions := putMessageOptions{
		StructID:   msg[376:380],
		Version:    msg[380:384],
		Options:    msg[384:388],
		Timeout:    msg[388:392],
		Context:    msg[392:396],
		KnDstCnt:   []byte{0x01, 0x00, 0x00, 0x00},
		UkDstCnt:   []byte{0x00, 0x00, 0x00, 0x00},
		InDstCnt:   []byte{0x00, 0x00, 0x00, 0x00},
		ResQName:   resQName,
		ResQMgr:    ctx.sessions[cid].qMgr,
		NumRecs:    msg[504:508],
		PMRFlag:    msg[508:512],
		OffsetPMR:  msg[512:516],
		OffsetRR:   msg[516:520],
		AddressPMR: msg[520:524],
		AddressRR:  msg[524:528],
	}
	response = append(response, getBytes(putMessageOptions)...)

	mqput := mqput{
		DataLength: msg[528:532],
	}
	response = append(response, getBytes(mqput)...)

	jmsValue, usrValue := getMsgInfo(msg)
	message := message{
		jmsValue: jmsValue,
		usrValue: usrValue,
	}
	q := ctx.sessions[cid].hdls[hdl].queue
	q.put(message)

	return response
}

func getQueueName(msg []byte) string {
	valueLen1 := int(binary.BigEndian.Uint32(msg[568:572]))
	valueLen2 := int(binary.BigEndian.Uint32(msg[572+valueLen1 : 572+valueLen1+4]))
	queueStr := string(msg[572+valueLen1+4 : 572+valueLen1+4+valueLen2])
	qRegEx := regexp.MustCompile(`///(.+?)<`)

	return qRegEx.FindStringSubmatch(queueStr)[1]
}

func getMsgInfo(msg []byte) (jmsValue, usrValue []byte) {
	valueLen1 := int(binary.BigEndian.Uint32(msg[568:572]))
	valueLen2 := int(binary.BigEndian.Uint32(msg[572+valueLen1 : 572+valueLen1+4]))
	valueLen3 := int(binary.BigEndian.Uint32(msg[572+valueLen1+4+valueLen2 : 572+valueLen1+4+valueLen2+4]))

	jmsValue = msg[572+valueLen1+4 : 572+valueLen1+4+valueLen2]
	usrValue = msg[572+valueLen1+4+valueLen2+4 : 572+valueLen1+4+valueLen2+4+valueLen3]

	return jmsValue, usrValue
}
