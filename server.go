package http

import (
	"net/http"
	"time"
)

type ServerConf struct {
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

var DefaultServer = &http.Server{
	ReadTimeout:       10 * time.Second,
	ReadHeaderTimeout: 5 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       10 * time.Minute,
}
