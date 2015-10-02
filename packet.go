package govertica

type SendablePacket interface {
	Payload() []byte
}
