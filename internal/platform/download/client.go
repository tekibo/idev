package download

import (
	"net"
	"net/http"
	"time"
)

const (
	maxAttempts = 8
	maxDelay    = 15 * time.Second
	maxFileSize = 8 << 30
)

var client = &http.Client{Transport: transport()}

func transport() *http.Transport {
	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	return &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   30 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		ForceAttemptHTTP2:     true,
	}
}
