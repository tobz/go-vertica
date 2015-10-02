package messages

import "bytes"
import "encoding/binary"

func integerAsPayload(i uint32) []byte {
	w := new(bytes.Buffer)
	binary.Write(w, binary.BigEndian, uint32(8))
	binary.Write(w, binary.BigEndian, i)

	return w.Bytes()
}
