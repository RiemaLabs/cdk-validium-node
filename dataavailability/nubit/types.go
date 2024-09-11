package nubit

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

const (
	NubitBatchDADataSize = 40
)

type BatchDAData struct {
	BlockNumber int64
	Commitment  [32]byte
}

// write a function that encode batchDAData struct into ABI-encoded bytes
func (b BatchDAData) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, b.BlockNumber)
	binary.Write(buf, binary.LittleEndian, b.Commitment)
	return buf.Bytes(), nil
}
func (b *BatchDAData) Decode(data []byte) error {
	buf := bytes.NewReader(data)
	binary.Read(buf, binary.LittleEndian, &b.BlockNumber)
	binary.Read(buf, binary.LittleEndian, &b.Commitment)
	return nil
}
func (b BatchDAData) IsEmpty() bool {
	return reflect.DeepEqual(b, BatchDAData{})
}
