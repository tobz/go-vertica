package govertica

import "fmt"
import "time"
import "strings"
import "errors"
import "net/url"
import "crypto/tls"

var (
	errInvalidDSNUnescaped       = errors.New("Invalid DSN: Did you forget to escape a param value?")
	errInvalidDSNAddr            = errors.New("Invalid DSN: Network Address not terminated (missing closing brace)")
	errInvalidDSNNoSlash         = errors.New("Invalid DSN: Missing the slash separating the database name")
	errInvalidDSNUnsafeCollation = errors.New("Invalid DSN: interpolateParams can be used with ascii, latin1, utf8 and utf8mb4 charset")
)

// parseDSN parses the DSN string to a config
func parseDSN(dsn string) (cfg *config, err error) {
	// New config with some default values
	cfg = &config{
		loc: time.UTC,
	}

	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	// Find the last '/' (since the password or the net addr might contain a '/')
	foundSlash := false
	for i := len(dsn) - 1; i >= 0; i-- {
		if dsn[i] == '/' {
			foundSlash = true
			var j, k int

			// left part is empty if i <= 0
			if i > 0 {
				// [username[:password]@][protocol[(address)]]
				// Find the last '@' in dsn[:i]
				for j = i; j >= 0; j-- {
					if dsn[j] == '@' {
						// username[:password]
						// Find the first ':' in dsn[:j]
						for k = 0; k < j; k++ {
							if dsn[k] == ':' {
								cfg.passwd = dsn[k+1 : j]
								break
							}
						}
						cfg.user = dsn[:k]

						break
					}
				}

				// [protocol[(address)]]
				// Find the first '(' in dsn[j+1:i]
				for k = j + 1; k < i; k++ {
					if dsn[k] == '(' {
						// dsn[i-1] must be == ')' if an address is specified
						if dsn[i-1] != ')' {
							if strings.ContainsRune(dsn[k+1:i], ')') {
								return nil, errInvalidDSNUnescaped
							}
							return nil, errInvalidDSNAddr
						}
						cfg.addr = dsn[k+1 : i-1]
						break
					}
				}
				cfg.net = dsn[j+1 : k]
			}

			// dbname[?param1=value1&...&paramN=valueN]
			// Find the first '?' in dsn[i+1:]
			for j = i + 1; j < len(dsn); j++ {
				if dsn[j] == '?' {
					if err = parseDSNParams(cfg, dsn[j+1:]); err != nil {
						return
					}
					break
				}
			}
			cfg.dbname = dsn[i+1 : j]

			break
		}
	}

	if !foundSlash && len(dsn) > 0 {
		return nil, errInvalidDSNNoSlash
	}

	// Set default network if empty
	if cfg.net == "" {
		cfg.net = "tcp"
	}

	// Set default address if empty
	if cfg.addr == "" {
		switch cfg.net {
		case "tcp":
			cfg.addr = "127.0.0.1:5534"
		default:
			return nil, errors.New("Default addr for network '" + cfg.net + "' unknown")
		}

	}

	return
}

// parseDSNParams parses the DSN "query string"
// Values must be url.QueryEscape'ed
func parseDSNParams(cfg *config, params string) (err error) {
	for _, v := range strings.Split(params, "&") {
		param := strings.SplitN(v, "=", 2)
		if len(param) != 2 {
			continue
		}

		// cfg params
		switch value := param[1]; param[0] {
		// Time Location
		case "loc":
			if value, err = url.QueryUnescape(value); err != nil {
				return
			}
			cfg.loc, err = time.LoadLocation(value)
			if err != nil {
				return
			}

			// Dial Timeout
		case "timeout":
			cfg.timeout, err = time.ParseDuration(value)
			if err != nil {
				return
			}

			// TLS-Encryption
		case "tls":
			boolValue, isBool := readBool(value)
			if isBool {
				if boolValue {
					cfg.tls = &tls.Config{}
				}
			} else {
				if strings.ToLower(value) == "skip-verify" {
					cfg.tls = &tls.Config{InsecureSkipVerify: true}
				} else {
					return fmt.Errorf("Invalid value / unknown config name: %s", value)
				}
			}

		default:
			// lazy init
			if cfg.params == nil {
				cfg.params = make(map[string]string)
			}

			if cfg.params[param[0]], err = url.QueryUnescape(value); err != nil {
				return
			}
		}
	}

	return
}

// Returns the bool value of the input.
// The 2nd return value indicates if the input was a valid bool value
func readBool(input string) (value bool, valid bool) {
	switch input {
	case "1", "true", "TRUE", "True":
		return true, true
	case "0", "false", "FALSE", "False":
		return false, true
	}

	// Not a valid bool value
	return
}
