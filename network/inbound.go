package network

import "fmt"
import "bytes"
import "github.com/tobz/go-vertica/constants"

type InboundPacket struct {
	Type constants.PacketType
	buf  *bytes.Buffer
}

func NewInboundPacket(buf []byte) *InboundPacket {
	packet := &InboundPacket{buf: bytes.NewBuffer(buf)}
	packetType := packet.ReadUInt8()
	packet.Type = constants.PacketType(packetType)

	// Read out a uint32 to skip the message size bytes.
	packet.ReadUInt32()

	return packet
}

func (self *InboundPacket) ReadUInt8() uint8 {
	buf, err := self.readBytes(1)
	if err != nil {
		return 0
	}

	return uint8(buf[0])
}

func (self *InboundPacket) ReadUInt16() uint16 {
	buf, err := self.readBytes(2)
	if err != nil {
		return 0
	}

	return uint16(buf[1]) | uint16(buf[0]<<8)
}

func (self *InboundPacket) ReadUInt32() uint32 {
	buf, err := self.readBytes(4)
	if err != nil {
		return 0
	}

	return uint32(buf[3]) | uint32(buf[2]<<8) | uint32(buf[1]<<16) | uint32(buf[0]<<24)
}

func (self *InboundPacket) ReadUInt64() uint64 {
	buf, err := self.readBytes(8)
	if err != nil {
		return 0
	}

	return uint64(buf[7] | buf[6]<<8 | buf[5]<<16 | buf[4]<<24 | buf[3]<<32 | buf[2]<<40 | buf[1]<<48 | buf[0]<<56)
}

func (self *InboundPacket) ReadRemaining() []byte {
	return self.buf.Next(self.buf.Len())
}

func (self *InboundPacket) ReadString() string {
	next, err := self.buf.ReadBytes('\x00')
	if err != nil {
		return ""
	}

	return string(next[0 : len(next)-1])
}

func (self *InboundPacket) Buffer() []byte {
	return self.buf.Bytes()
}

func (self *InboundPacket) readBytes(count int) ([]byte, error) {
	if self.buf.Len() < count {
		return nil, fmt.Errorf("tried to read %d bytes from packet, only %d left! (id: %d)", count, self.buf.Len(), uint32(self.Type))
	}

	return self.buf.Next(count), nil
}
