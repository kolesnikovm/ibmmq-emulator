package mq

type socketAction struct {
	ConversID []byte
	RequestID []byte
	Type      []byte
	Param1    []byte
	Param2    []byte
}
