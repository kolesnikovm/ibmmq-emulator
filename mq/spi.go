package mq

import (
	"bytes"
	"encoding/hex"
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
	QUERY   = []byte{0x01, 0x00, 0x00, 0x00}
	OPEN    = []byte{0x0c, 0x00, 0x00, 0x00}
)

func handleSpi(msg []byte) (response []byte) {
	spiVerb := msg[:4]

	spi := spi{
		SpiVerb:      spiVerb,
		Version:      msg[4:8],
		MaxReplySize: msg[8:12],
	}
	response = append(response, spi.SpiVerb...)
	response = append(response, spi.Version...)
	response = append(response, spi.MaxReplySize...)

	switch {
	case bytes.Compare(spiVerb, QUERY) == 0:
		spqu := spib{
			SpiStructID: msg[12:16],
			Version:     msg[16:20],
			Length:      msg[20:24],
		}
		response = append(response, spqu.SpiStructID...)
		response = append(response, spqu.Version...)
		response = append(response, spqu.Length...)

		response = append(response, spqo...)
	case bytes.Compare(spiVerb, OPEN) == 0:
		spou := spib{
			SpiStructID: msg[12:16],
			Version:     msg[16:20],
			Length:      msg[20:24],
		}
		response = append(response, getBytes(spou)...)

		spoo := spib{
			SpiStructID: []byte{0x53, 0x50, 0x4f, 0x4f},
			Version:     msg[28:32],
			Length:      []byte{0x04, 0x02, 0x00, 0x00},
		}
		response = append(response, getBytes(spoo)...)

		qproject, _ := hex.DecodeString("53595354454d2e50524f54454354494f4e2e4552524f522e515545554520202020202020202020202020202020202020")
		lpoo := lpoo{
			StructID:       msg[36:40],
			Version:        msg[40:44],
			LpiVersion:     msg[44:48],
			LpiOpts:        msg[48:52],
			DefPersistence: []byte{0x00, 0x00, 0x00, 0x00},
			DefPutRespType: []byte{0x01, 0x00, 0x00, 0x00},
			DefReadAhead:   []byte{0x00, 0x00, 0x00, 0x00},
			PropertyCtl:    []byte{0x00, 0x00, 0x00, 0x00},
			QProtect:       qproject,
			Options:        []byte{0x2c, 0x01, 0x00, 0x00},
			ExtraData:      []byte{0xd8, 0x01, 0x00, 0x00},
		}
		response = append(response, getBytes(lpoo)...)

		objectDescriptor := objectDescriptor{
			StructID:       msg[124:128],
			Version:        msg[128:132],
			ObjType:        msg[132:136],
			ObjName:        msg[136:184],
			ObjQMgr:        msg[184:232],
			DynQName:       msg[232:280],
			AltUserID:      msg[280:292],
			NbrRecord:      msg[292:296],
			KnownDestCount: []byte{0x01, 0x00, 0x00, 0x00},
			UnknownDestCnt: msg[300:304],
			InvalidDestCnt: msg[304:308],
			OffsetOR:       msg[308:312],
			OffsetRR:       msg[312:316],
			AddressOR:      msg[316:320],
			AddressRR:      msg[320:324],
		}
		response = append(response, getBytes(objectDescriptor)...)

		XtradStart, _ := hex.DecodeString("000000080000002400000003000000000000000100000001000000000000000000000008000000030000001000000005000000010000000400000044000000060000033300000030")
		XtradEnd, _ := hex.DecodeString("00000003000000100000000700000000000000030000001000000008000000000000000300000010000000090000000000000003000000100000000e0000000000000003000000100000000a0000000000000003000000100000000c00000000")
		response = append(response, XtradStart...)
		response = append(response, msg[136:184]...)
		response = append(response, XtradEnd...)
	default:
		log.Println("[WARN] unknown spi verb")
	}

	return response
}
