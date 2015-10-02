package network

import "bytes"
import "github.com/tobz/go-vertica/constants"

type OutboundPacket struct {
	packetType constants.PacketType
	buf        bytes.Buffer
	finalized  bool
}

func NewOutboundPacket(packetType constants.PacketType) *OutboundPacket {
	packet := &OutboundPacket{packetType: packetType}
	packet.WriteUInt8(uint8(packetType))

	// Message size placeholder.
	packet.WriteUInt32(0)

	return packet
}

func NewOutboundPacketWithoutType() *OutboundPacket {
	packet := &OutboundPacket{}

	// Message size placeholder.
	packet.WriteUInt32(0)

	return packet
}

func (self *OutboundPacket) WriteUInt8(val uint8) {
	self.buf.Write([]byte{byte(val)})
}

func (self *OutboundPacket) WriteUInt16(val uint16) {
	self.buf.Write([]byte{byte(val >> 8), byte(val)})
}

func (self *OutboundPacket) WriteUInt32(val uint32) {
	self.buf.Write([]byte{byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)})
}

func (self *OutboundPacket) WriteUInt64(val uint64) {
	self.buf.Write([]byte{byte(val >> 56), byte(val >> 48), byte(val >> 40), byte(val >> 32), byte(val >> 24), byte(val >> 16), byte(val >> 8), byte(val)})
}

func (self *OutboundPacket) WriteStringTuple(key, value string) {
	self.WriteString(key)
	self.WriteString(value)
}

func (self *OutboundPacket) WriteString(value string) {
	self.buf.WriteString(value)
	self.buf.WriteByte(0x00)
}

func (self *OutboundPacket) WriteNull() {
	self.buf.WriteByte(0x00)
}

func (self *OutboundPacket) finalize() {
	if self.finalized {
		return
	}

	buf := self.buf.Bytes()
	bufLen := len(buf) - 1
	offset := 1

	if self.packetType == constants.PacketTypeNone {
		bufLen = bufLen + 1
		offset = 0
	}

	buf[offset+0] = byte(bufLen >> 24)
	buf[offset+1] = byte(bufLen >> 16)
	buf[offset+2] = byte(bufLen >> 8)
	buf[offset+3] = byte(bufLen)

	self.finalized = true
}

func (self *OutboundPacket) Construct() []byte {
	self.finalize()
	return self.buf.Bytes()
}
