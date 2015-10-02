package messages

import "github.com/tobz/go-vertica/network"

type BackendKeyResponse struct {
	Pid uint32
	Key uint32
}

func NewBackendKeyResponse(packet *network.InboundPacket) *BackendKeyResponse {
	pid := packet.ReadUInt32()
	key := packet.ReadUInt32()

	return &BackendKeyResponse{pid, key}
}
