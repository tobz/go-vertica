package messages

import "github.com/tobz/go-vertica/constants"
import "github.com/tobz/go-vertica/network"

type AuthenticationResponse struct {
	Result   constants.AuthenticationResult
	Salt     []byte
	AuthData []byte
}

func NewAuthenticationResponse(packet *network.InboundPacket) *AuthenticationResponse {
	var salt []byte = nil
	var authData []byte = nil

	resultRaw := packet.ReadUInt32()
	result := constants.AuthenticationResult(resultRaw)
	other := packet.ReadRemaining()

	if result == constants.AuthenticationResultCryptPassword || result == constants.AuthenticationResultMD5Password {
		salt = other
	}

	if result == constants.AuthenticationResultGSSContinue {
		authData = other
	}

	return &AuthenticationResponse{result, salt, authData}
}
