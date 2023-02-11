package models

import "errors"

var (
	ErrShutdownTimeout = errors.New("shutdown timeout")
	ErrStartTimeout    = errors.New("start timeout")
)
