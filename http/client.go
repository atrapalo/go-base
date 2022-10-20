package http

import (
	"net/http"
	"time"
)

const defaultTimeout = time.Second * 5

type ClientHttpCore struct {
	http.Client
}

func NewHttpCore() ClientHttpCore {
	return ClientHttpCore{
		Client: http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (c *ClientHttpCore) Do(req *http.Request) (*http.Response, error) {
	return c.Client.Do(req)
}

func (c *ClientHttpCore) Timeout(timeout time.Duration) {
	c.Client.Timeout = timeout
}
