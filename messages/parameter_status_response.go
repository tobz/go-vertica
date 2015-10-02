package messages

import "github.com/tobz/go-vertica/network"

type ParameterStatusResponse struct {
	Key   string
	Value string
}

func NewParameterStatusResponse(packet *network.InboundPacket) *ParameterStatusResponse {
	key := packet.ReadString()
	value := packet.ReadString()

	return &ParameterStatusResponse{key, value}
}
