package constants

type MessageType rune

const MessageTypeNone MessageType = '0'
const MessageTypeAuthenticationResponse MessageType = 'R'
const MessageTypeReadyToQueryResponse MessageType = 'Z'

const AuthenticationResponseOK int = 0
