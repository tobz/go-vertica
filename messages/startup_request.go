package messages

import "github.com/tobz/go-vertica/constants"
import "github.com/tobz/go-vertica/network"

type StartupRequest struct {
	username string
	database string
	options  string
}

func NewStartupRequest(username, database, options string) *StartupRequest {
	return &StartupRequest{username, database, options}
}

func (self *StartupRequest) Type() constants.MessageType {
	return constants.MessageTypeNone
}

func (self *StartupRequest) Payload() []byte {
	packet := network.NewOutboundPacketWithoutType()
	packet.WriteUInt32(uint32(constants.VerticaProtocolVersion))

	if self.username != "" {
		packet.WriteStringTuple("user", self.username)
	}

	if self.database != "" {
		packet.WriteStringTuple("database", self.database)
	}

	if self.options != "" {
		packet.WriteStringTuple("options", self.options)
	}

	packet.WriteNull()

	return packet.Construct()
}
