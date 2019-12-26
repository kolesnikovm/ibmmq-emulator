package mq

type notification struct {
	version []byte
	handle  []byte
	code    []byte
	value   []byte
}

var (
	CHECK_MSG = []byte{0x0b, 0x00, 0x00, 0x00}
)

func handleNotification(msg []byte) (response []byte) {
	tshm := tshm{
		StructID:  []byte(TSHM),
		MQSegmLen: []byte{0x00, 0x00, 0x00, 0x34},
		ConversID: msg[8:12],
		RequestID: msg[12:16],
		ByteOrder: msg[16:17],
		SegmType:  []byte{NOTIFICATION},
		CtlFlag1:  []byte{0x10},
		CtlFlag2:  msg[19:20],
		LUWIdent:  msg[20:28],
		Encoding:  REVERSED,
		CCSID:     msg[32:34],
		Reserved:  msg[34:36],
	}
	response = append(response, getBytes(tshm)...)

	notification := notification{
		version: []byte{0x01, 0x00, 0x00, 0x00},
		handle:  msg[40:44],
		code:    CHECK_MSG,
		value:   []byte{0x00, 0x00, 0x00, 0x00},
	}
	response = append(response, getBytes(notification)...)

	return response
}
