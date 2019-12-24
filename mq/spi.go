package mq

import (
	"bytes"
	"encoding/hex"
	"encoding/binary"
	"log"
)

type spi struct {
	SpiVerb      []byte `offset:"0", length:"4"`
	Version      []byte `offset:"4", length:"4"`
	MaxReplySize []byte `offset:"8", length:"4"`
}

type spib struct {
	SpiStructID []byte `offset:"0", length:"4"`
	Version     []byte `offset:"4", length:"4"`
	Length      []byte `offset:"8", length:"4"`
}

type lpoo struct {
	StructID       []byte
	Version        []byte
	LpiVersion     []byte
	LpiOpts        []byte
	DefPersistence []byte
	DefPutRespType []byte
	DefReadAhead   []byte
	PropertyCtl    []byte
	QProtect       []byte
	Options        []byte
	ExtraData      []byte
}

var (
	spqo, _ = hex.DecodeString("5350514f01000000000100000c00000001000000020000000200000002000000010000000200000001000000030000000100000001000000030000000200000002000000030000000100000004000000010000000100000001000000010000000500000001000000010000000100000001000000060000000100000001000000010000000100000007000000010000000100000001000000010000000800000001000000010000000100000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000b000000010000000100000001000000010000000c00000001000000020000000200000001000000")

	QUERY = []byte{0x01, 0x00, 0x00, 0x00}
	OPEN  = []byte{0x0c, 0x00, 0x00, 0x00}

	SPQU = []byte{0x53, 0x50, 0x51, 0x55}
	SPQO = []byte{0x53, 0x50, 0x51, 0x4f}

	SPOU = []byte{0x53, 0x50, 0x4f, 0x55}
	SPOO = []byte{0x53, 0x50, 0x4f, 0x4f}
	LPOO = []byte{0x4c, 0x50, 0x4f, 0x4f}
	OD   = []byte{0x4f, 0x44, 0x20, 0x20}
)

func handleSpi(msg []byte) (response []byte) {
	spiVerb := msg[52:56]

	spi := spi{
		SpiVerb:      spiVerb,
		Version:      msg[56:60],
		MaxReplySize: msg[60:64],
	}
	response = append(response, spi.SpiVerb...)
	response = append(response, spi.Version...)
	response = append(response, spi.MaxReplySize...)

	switch {
	case bytes.Compare(spiVerb, QUERY) == 0:
		log.Printf("[INFO] M: SPI, C: %d, R: %d, V: QUERY\n", binary.BigEndian.Uint32(msg[8:12]), binary.BigEndian.Uint32(msg[12:16]))

		spqu := spib{
			SpiStructID: SPQU,
			Version:     msg[68:72],
			Length:      msg[72:76],
		}
		response = append(response, getBytes(spqu)...)
		response = append(response, spqo...)
	case bytes.Compare(spiVerb, OPEN) == 0:
		log.Printf("[INFO] M: SPI, C: %d, R: %d, V: OPEN, Obj: %s\n", binary.BigEndian.Uint32(msg[8:12]), binary.BigEndian.Uint32(msg[12:16]), msg[188:236])

		spou := spib{
			SpiStructID: SPOU,
			Version:     msg[68:72],
			Length:      msg[72:76],
		}
		response = append(response, getBytes(spou)...)

		spoo := spib{
			SpiStructID: SPOO,
			Version:     msg[80:84],
			Length:      []byte{0x04, 0x02, 0x00, 0x00},
		}
		response = append(response, getBytes(spoo)...)

		lpoo := lpoo{
			StructID:       LPOO,
			Version:        msg[92:96],
			LpiVersion:     msg[96:100],
			LpiOpts:        msg[100:104],
			DefPersistence: []byte{0x00, 0x00, 0x00, 0x00},
			DefPutRespType: []byte{0x01, 0x00, 0x00, 0x00},
			DefReadAhead:   []byte{0x00, 0x00, 0x00, 0x00},
			PropertyCtl:    []byte{0x00, 0x00, 0x00, 0x00},
			QProtect:       decodeString("53595354454d2e50524f54454354494f4e2e4552524f522e515545554520202020202020202020202020202020202020"),
			Options:        []byte{0x2c, 0x01, 0x00, 0x00},
			ExtraData:      []byte{0xd8, 0x01, 0x00, 0x00},
		}
		response = append(response, getBytes(lpoo)...)

		objectDescriptor := objectDescriptor{
			StructID:       OD,
			Version:        msg[180:184],
			ObjType:        msg[184:188],
			ObjName:        msg[188:236],
			ObjQMgr:        msg[236:284],
			DynQName:       msg[284:332],
			AltUserID:      msg[332:344],
			NbrRecord:      msg[344:348],
			KnownDestCount: []byte{0x01, 0x00, 0x00, 0x00},
			UnknownDestCnt: msg[352:356],
			InvalidDestCnt: msg[356:360],
			OffsetOR:       msg[360:364],
			OffsetRR:       msg[364:368],
			AddressOR:      msg[368:372],
			AddressRR:      msg[372:376],
		}
		response = append(response, getBytes(objectDescriptor)...)

		XtradStart, _ := hex.DecodeString("000000080000002400000003000000000000000100000001000000000000000000000008000000030000001000000005000000010000000400000044000000060000033300000030")
		XtradEnd, _ := hex.DecodeString("00000003000000100000000700000000000000030000001000000008000000000000000300000010000000090000000000000003000000100000000e0000000000000003000000100000000a0000000000000003000000100000000c00000000")
		response = append(response, XtradStart...)
		response = append(response, msg[188:236]...)
		response = append(response, XtradEnd...)
	default:
		log.Printf("[WARN] unknown spi verb: %x\n", spiVerb)
	}

	return response
}
