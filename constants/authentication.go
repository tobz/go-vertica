package constants

type AuthenticationResult uint32

const (
	AuthenticationResultOK                AuthenticationResult = 0
	AuthenticationResultKerberosV5        AuthenticationResult = 2
	AuthenticationResultCleartextPassword AuthenticationResult = 3
	AuthenticationResultCryptPassword     AuthenticationResult = 4
	AuthenticationResultMD5Password       AuthenticationResult = 5
	AuthenticationResultSCMCredential     AuthenticationResult = 6
	AuthenticationResultGSS               AuthenticationResult = 7
	AuthenticationResultGSSContinue       AuthenticationResult = 8
	AuthenticationResultSSPI              AuthenticationResult = 9
)
