package constants

type PacketType uint8

const (
	PacketTypeNone                    PacketType = 0
	PacketTypeAuthenticationResponse  PacketType = 'R'
	PacketTypeReadyForQueryResponse   PacketType = 'Z'
	PacketTypePasswordRequest         PacketType = 'p'
	PacketTypeServerNoticeResponse    PacketType = 'N'
	PacketTypeServerErrorResponse     PacketType = 'E'
	PacketTypeParameterStatusResponse PacketType = 'S'
	PacketTypeBackendKeyResponse      PacketType = 'K'
)
