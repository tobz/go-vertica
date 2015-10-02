package messages

import "github.com/tobz/go-vertica/constants"
import "github.com/tobz/go-vertica/network"

var messageFieldMappings = map[string]string{
	"q": "Internal Query",
	"S": "Severity",
	"M": "Message",
	"C": "Sqlstate",
	"D": "Detail",
	"H": "Hint",
	"P": "Position",
	"W": "Where",
	"p": "Internal Position",
	"R": "Routine",
	"F": "File",
	"L": "Line",
}

type ServerMessageResponse struct {
	Type   string
	Fields map[string]string
}

func NewServerMessageResponse(packet *network.InboundPacket) *ServerMessageResponse {
	messageType := "notice"

	if packet.Type == constants.PacketTypeServerErrorResponse {
		messageType = "error"
	}

	messageFields := make(map[string]string)

	for {
		field := packet.ReadString()
		if field == "" {
			break
		}

		fieldType := string([]byte(field)[0])
		fieldValue := string([]byte(field)[1:])

		messageFields[messageFieldMappings[fieldType]] = fieldValue
	}

	return &ServerMessageResponse{messageType, messageFields}
}
