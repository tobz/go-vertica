package network

import "io"
import "log"

type PacketReader struct {
	conn    io.Reader
	readBuf []byte
	readOff int
}

func NewPacketReader(conn io.Reader, bufSize int) *PacketReader {
	return &PacketReader{conn, make([]byte, bufSize), 0}
}

func (r *PacketReader) Next() (*InboundPacket, error) {
	// See if we have any packets in our left over buffer.
	bufferedPacket, err := r.getPacketFromBuffer()
	if bufferedPacket != nil || err != nil {
		return bufferedPacket, err
	}

	// We've cleared through any buffered packets, so do a network read to get everything off the wire.
	n, err := r.conn.Read(r.readBuf[r.readOff:])
	if err != nil {
		return nil, err
	}

	// If we didn't get anything, then we won't be able to get a packet out.
	if n == 0 {
		return nil, nil
	}

	// Increment our read offset based on what we received.
	r.readOff += n

	// See if we can get a packet out yet.
	return r.getPacketFromBuffer()
}

func (r *PacketReader) getPacketFromBuffer() (*InboundPacket, error) {
	// Make sure we have a packet in the buffer to read.
	if !r.hasBufferedPacket() {
		return nil, nil
	}

	packetLength := r.getPacketLengthFromBuffer()

	log.Printf("reader: received packet -> %#v", r.readBuf[:packetLength])

	buf := make([]byte, packetLength)
	copy(buf, r.readBuf[:packetLength])

	packet := NewInboundPacket(buf)

	// If we pulled out a packet and there's left over data, we need to shift it to the front so the next call
	// will read it out immediately.  Otherwise, we're done here, so put the read offset back to zero to start
	// reading into our buffer from the front.
	if r.readOff > packetLength {
		remaining := (r.readOff - packetLength)
		copy(r.readBuf, r.readBuf[packetLength:])
		r.readOff = remaining
	} else {
		r.readOff = 0
	}

	return packet, nil
}

func (r *PacketReader) hasBufferedPacket() bool {
	// See if we have enough bytes in our buffer for a packet at all.
	if r.readOff < 5 {
		return false
	}

	packetLength := r.getPacketLengthFromBuffer()

	if r.readOff < packetLength {
		return false
	}

	return true
}

func (r *PacketReader) getPacketLengthFromBuffer() int {
	// Incoming packets have [id - 1 byte][message size - 4 bytes][message - n bytes].  Message size
	// is inclusive of itself, so we just need to account for the message ID which, for some reason,
	// comes before the message size... and message size is in network order.
	return (int(r.readBuf[4]) | int(r.readBuf[3]<<8) | int(r.readBuf[2]<<16) | int(r.readBuf[1]<<24)) + 1
}

type PacketWriter struct {
	conn io.Writer
}

func NewPacketWriter(conn io.Writer) *PacketWriter {
	return &PacketWriter{conn}
}

func (w *PacketWriter) Write(buf []byte) error {
	_, err := w.conn.Write(buf)

	log.Printf("writer: sending packet -> %#v", buf)

	return err
}
