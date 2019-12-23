package mq

import (
	"encoding/hex"
	"log"
)

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

type mqOpen struct {
	Options []byte `offset:"0", length:"4"`
}

type fopa struct {
	StructID        []byte `offset:"0", length:"4"`
	Version         []byte `offset:"4", length:"4"`
	Length          []byte `offset:"8", length:"4"`
	DefPersistence  []byte `offset:"12", length:"4"`
	DefPutRespType  []byte `offset:"16", length:"4"`
	DefReadAhead    []byte `offset:"20", length:"4"`
	PropertyControl []byte `offset:"24", length:"4"`
	hiddenField     []byte `offset:"28", length:"56"`
}

func handleMqOpen(msg []byte) (response []byte) {
	log.Printf("[INFO] received MQOPEN message\n")

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
	response = append(response, getBytes(objectDescriptor)...)

	mqOpen := mqOpen{
		Options: []byte{0x20, 0x20, 0x00, 0x00},
	}
	response = append(response, getBytes(mqOpen)...)

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
	response = append(response, getBytes(fopa)...)

	return response
}
