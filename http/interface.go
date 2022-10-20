package http

import (
	"net/http"
	"time"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
	Timeout(timeout time.Duration)
}
