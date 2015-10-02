package interfaces

import "github.com/tobz/go-vertica/constants"

type Message interface {
	Type() constants.MessageType
	Payload() []byte
}
