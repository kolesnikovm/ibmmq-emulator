package mq

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
