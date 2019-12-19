package mq

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type mqInc struct {
	SelectorCount   []byte     `offset:"0" length:"4"`
	IntegerCount    []byte     `offset:"4" length:"4"`
	CharacterLength []byte     `offset:"8" length:"4"`
	Selectors       []selector `offset:"12"`
	IntegerValues   []integerValue
	CharValues      []byte
}

type selector struct {
	value []byte `length:"4"`
}

type integerValue struct {
	value []byte `length:"4"`
}

func handleMqInc(msg []byte) (byteValue []byte) {
	selectorCount := int(binary.LittleEndian.Uint32(msg[:4]))
	integerCount := int(binary.LittleEndian.Uint32(msg[4:8]))
	characterCount := int(binary.LittleEndian.Uint32(msg[8:12])) / 48

	selectors := make([]selector, selectorCount)
	for i := 0; i < selectorCount; i++ {
		selectors[i] = selector{
			value: msg[12+i*4 : 16+i*4],
		}
	}

	integerValues := make([]integerValue, integerCount)
	for i := 0; i < integerCount; i++ {
		integerValues[i] = integerValue{
			value: getSelectorValue(selectors[i].value),
		}
	}

	characterValues := make([]byte, 0, characterCount*48)
	for i := 0; i < characterCount; i++ {
		characterValues = append(characterValues, getSelectorValue(selectors[integerCount+i].value)...)
	}

	mqInc := mqInc{
		SelectorCount:   msg[0:4],
		IntegerCount:    msg[4:8],
		CharacterLength: msg[8:12],
		Selectors:       selectors,
		IntegerValues:   integerValues,
		CharValues:      characterValues,
	}

	byteValue = append(byteValue, mqInc.SelectorCount...)
	byteValue = append(byteValue, mqInc.IntegerCount...)
	byteValue = append(byteValue, mqInc.CharacterLength...)
	for _, selector := range mqInc.Selectors {
		byteValue = append(byteValue, selector.value...)
	}
	for _, intValue := range mqInc.IntegerValues {
		byteValue = append(byteValue, intValue.value...)
	}
	byteValue = append(byteValue, mqInc.CharValues...)

	return byteValue
}

func getSelectorValue(key []byte) (value []byte) {
	switch {
	case bytes.Compare(key, []byte{0x1f, 0x00, 0x00, 0x00}) == 0:
		value = []byte{0x90, 0x03, 0x00, 0x00}
	case bytes.Compare(key, []byte{0x20, 0x00, 0x00, 0x00}) == 0:
		value = []byte{0x03, 0x00, 0x00, 0x00}
	case bytes.Compare(key, []byte{0x02, 0x00, 0x00, 0x00}) == 0:
		value = []byte{0x33, 0x03, 0x00, 0x00}
	case bytes.Compare(key, []byte{0xdf, 0x07, 0x00, 0x00}) == 0:
		value, _ = hex.DecodeString("514d31202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020")
		// queue manager name
	case bytes.Compare(key, []byte{0xf0, 0x07, 0x00, 0x00}) == 0:
		value, _ = hex.DecodeString("514d315f323031392d31322d30355f31312e31372e303120202020202020202020202020202020202020202020202020")
		// queue manager identifier
	case bytes.Compare(key, []byte{0xd6, 0x07, 0x00, 0x00}) == 0:
		value, _ = hex.DecodeString("4445562e444541442e4c45545445522e5155455545202020202020202020202020202020202020202020202020202020")
		// dead letter queue
	}

	return value
}
