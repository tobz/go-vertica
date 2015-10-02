package messages

import "fmt"
import "crypto/md5"
import "github.com/tobz/go-vertica/constants"
import "github.com/tobz/go-vertica/network"

type PasswordRequest struct {
	username string
	password string
	authResp *AuthenticationResponse
}

func NewPasswordRequest(username, password string, authResp *AuthenticationResponse) *PasswordRequest {
	return &PasswordRequest{username, password, authResp}
}

func (self *PasswordRequest) generatePassword() string {
	switch self.authResp.Result {
	case constants.AuthenticationResultCleartextPassword:
		return self.password
	case constants.AuthenticationResultCryptPassword:
		return cryptPassword(self.password, self.authResp.Salt)
	case constants.AuthenticationResultMD5Password:
		return md5Password(self.username, self.password, self.authResp.Salt)
	default:
		panic(fmt.Sprintf("authentication result type %d not supported!", uint32(self.authResp.Result)))
	}
}

func (self *PasswordRequest) Payload() []byte {
	packet := network.NewOutboundPacket(constants.PacketTypePasswordRequest)

	password := self.generatePassword()
	packet.WriteString(password)

	return packet.Construct()
}

func cryptPassword(password string, salt []byte) string {
	return ""
}

func md5Password(username, password string, salt []byte) string {
	firstPass := fmt.Sprintf("%x", md5.Sum([]byte(password+username)))
	secondPass := fmt.Sprintf("%x", md5.Sum(append([]byte(firstPass), salt...)))
	return fmt.Sprintf("md5%s", secondPass)
}
