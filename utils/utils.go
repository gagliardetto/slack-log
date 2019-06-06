package utils

import (
	"net"
	"net/http"
	"time"
)

// default values
var (
	DefaultMaxIdleConnsPerHost               = 1000
	DefaultTimeout             time.Duration = 15 * time.Second
	DefaultKeepAlive           time.Duration = 180 * time.Second
)

// NewHTTPClient returns a new Client from the provided config.
// Client is safe for concurrent use by multiple goroutines.
func NewHTTPClient() *http.Client {

	tr := &http.Transport{
		//TLSClientConfig: &tls.Config{
		//	InsecureSkipVerify: conf.InsecureSkipVerify,
		//},
		MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
		Proxy:               http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   DefaultTimeout,
			KeepAlive: DefaultKeepAlive,
		}).Dial,
		TLSHandshakeTimeout: DefaultTimeout,
	}

	return &http.Client{
		Timeout:   DefaultTimeout,
		Transport: tr,
	}
}
