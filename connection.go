package govertica

import "time"
import "crypto/tls"

type config struct {
	user    string
	passwd  string
	net     string
	addr    string
	dbname  string
	params  map[string]string
	loc     *time.Location
	tls     *tls.Config
	timeout time.Duration
}
