package govertica

import "log"
import "time"
import "net"
import "crypto/tls"
import "github.com/tobz/go-vertica/constants"
import "github.com/tobz/go-vertica/network"
import "github.com/tobz/go-vertica/messages"

type Configuration struct {
	Network        string
	Address        string
	Username       string
	Password       string
	Database       string
	Params         map[string]string
	Location       *time.Location
	TLS            *tls.Config
	NetworkTimeout time.Duration
}

type Connection struct {
	config *Configuration
	conn   net.Conn
	isSsl  bool

	netRead     *network.PacketReader
	netWrite    *network.PacketWriter
	packetQueue chan *network.InboundPacket

	sessionId  int64
	parameters map[string]string
	backendPid uint32
	backendKey uint32
}

func NewConnection(dsn string) (*Connection, error) {
	config, err := parseDSN(dsn)
	if err != nil {
		return nil, err
	}

	return &Connection{config: config, packetQueue: make(chan *network.InboundPacket), parameters: make(map[string]string)}, nil
}

func (self *Connection) Connect() error {
	log.Printf("dialing %s...", self.config.Address)

	conn, err := net.DialTimeout(self.config.Network, self.config.Address, self.config.NetworkTimeout)
	if err != nil {
		return err
	}

	self.conn = conn
	self.netRead = network.NewPacketReader(conn, 32768)
	self.netWrite = network.NewPacketWriter(conn)

	log.Printf("connected! starting network loop...")

	go self.processNetwork()

	return self.authenticate()
}

func (self *Connection) authenticate() error {
	errors := make(chan error, 1)

	log.Printf("starting authentication loop...")

	self.takeNetworkLoop(func(packets <-chan *network.InboundPacket) {
		log.Printf("authenticate: sending startup request: user = %s, database = %s", self.config.Username, self.config.Database)
		startup := messages.NewStartupRequest(self.config.Username, self.config.Database, "")
		self.Send(startup)

		for {
			packet := <-packets
			switch packet.Type {
			case constants.PacketTypeAuthenticationResponse:
				authentication := messages.NewAuthenticationResponse(packet)
				log.Printf("authenticate: got response (code: %d)", uint32(authentication.Result))

				if authentication.Result != constants.AuthenticationResultOK {
					log.Printf("authenticate: sending password '%s' for username '%s'", self.config.Password, self.config.Username)
					password := messages.NewPasswordRequest(self.config.Username, self.config.Password, authentication)
					self.Send(password)
				}
			case constants.PacketTypeReadyForQueryResponse:
				log.Printf("authenticate: ready for query")
				// We're good to go, so exit the loop.
				return
			default:
				self.ProcessGenericMessage(packet)
			}
		}
	})

	err, ok := <-errors
	if ok {
		return err
	}

	return nil
}

func (self *Connection) handleError(err error) {
}

func (self *Connection) processNetwork() {
	for {
		p, err := self.netRead.Next()
		if err != nil {
			self.handleError(err)
			return
		}

		if p != nil {
			self.packetQueue <- p
		}
	}
}

func (self *Connection) takeNetworkLoop(handler func(packets <-chan *network.InboundPacket)) {
	loopFinished := make(chan struct{})

	go func() {
		// Handle the messages until the loop handler returns control.
		handler(self.packetQueue)

		// Tell our parent caller we're done.
		loopFinished <- struct{}{}
	}()

	// Wait for our loop handler to finish up.
	<-loopFinished
}

func (self *Connection) ProcessGenericMessage(packet *network.InboundPacket) {
	log.Printf("got generic packet '%d'...", uint8(packet.Type))

	switch packet.Type {
	case constants.PacketTypeServerNoticeResponse, constants.PacketTypeServerErrorResponse:
		message := messages.NewServerMessageResponse(packet)
		log.Printf("general: server message: %#v", message)
	case constants.PacketTypeParameterStatusResponse:
		message := messages.NewParameterStatusResponse(packet)
		self.parameters[message.Key] = message.Value

		log.Printf("general: parameters: %#v", self.parameters)
	case constants.PacketTypeBackendKeyResponse:
		message := messages.NewBackendKeyResponse(packet)
		self.backendPid = message.Pid
		self.backendKey = message.Key

		log.Printf("general: backend -> key: %d, pid: %d", self.backendKey, self.backendPid)
	}
}

func (self *Connection) Send(packet SendablePacket) {
	self.netWrite.Write(packet.Payload())
}
