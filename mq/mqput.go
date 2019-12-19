package mq

import (
	"encoding/binary"
	"encoding/hex"
	"regexp"
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

func handleMqPut(msg, userID, appType, appName, qMgr []byte) (response []byte) {
	msgID, _ := hex.DecodeString("414d5120514d312020202020202020206445ea5d04b59424") //сделать вручную
	messageDescriptor := messageDescriptor{
		StructID:   msg[0:4],
		Version:    msg[4:8],
		Report:     msg[8:12],
		MsgType:    msg[12:16],
		Expiry:     msg[16:20],
		Feedback:   msg[20:24],
		Encoding:   msg[24:28],
		CCSID:      msg[28:32],
		Format:     msg[32:40],
		Priority:   msg[40:44],
		Persist:    msg[44:48],
		MsgID:      msgID,
		CorrelID:   msg[72:96],
		BackoCnt:   msg[96:100],
		ReplyToQ:   msg[100:148],
		ReplToQMgr: msg[148:196],
		UserID:     userID,
		AccntTok:   msg[208:240],
		AppIDDAta:  msg[240:272],
		PutAppTyp:  appType,
		PutAppName: appName,
		PutDatGMT:  []byte{0x32, 0x30, 0x31, 0x39, 0x31, 0x32, 0x30, 0x36}, //сделать норм
		PutTimGMT:  []byte{0x31, 0x32, 0x34, 0x34, 0x31, 0x33, 0x35, 0x37},
		AppOriDat:  []byte{0x20, 0x20, 0x20, 0x20},
	}
	response = append(response, getBytes(messageDescriptor)...)

	valueLen1 := int(binary.BigEndian.Uint32(msg[516:520]))
	valueLen2 := int(binary.BigEndian.Uint32(msg[520+valueLen1 : 520+valueLen1+4]))
	queueStr := string(msg[520+valueLen1+4 : 520+valueLen1+4+valueLen2])
	qRegEx := regexp.MustCompile(`///(.+?)<`)
	queue := qRegEx.FindStringSubmatch(queueStr)[1]

	resQName := make([]byte, 0, 48)
	resQName = append(resQName, []byte(queue)...)
	for i := len(resQName); i < cap(resQName); i++ {
		resQName = append(resQName, []byte(" ")...)
	}

	putMessageOptions := putMessageOptions{
		StructID:   msg[324:328],
		Version:    msg[328:332],
		Options:    msg[332:336],
		Timeout:    msg[336:340],
		Context:    msg[340:344],
		KnDstCnt:   []byte{0x01, 0x00, 0x00, 0x00},
		UkDstCnt:   []byte{0x00, 0x00, 0x00, 0x00},
		InDstCnt:   []byte{0x00, 0x00, 0x00, 0x00},
		ResQName:   resQName,
		ResQMgr:    qMgr,
		NumRecs:    msg[452:456],
		PMRFlag:    msg[456:460],
		OffsetPMR:  msg[460:464],
		OffsetRR:   msg[464:468],
		AddressPMR: msg[468:472],
		AddressRR:  msg[472:476],
	}
	response = append(response, getBytes(putMessageOptions)...)

	mqput := mqput{
		DataLength: msg[476:480],
	}
	response = append(response, getBytes(mqput)...)

	return response
}
