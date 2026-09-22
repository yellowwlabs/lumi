package node

import "errors"

var (
	ErrNodeNotFound          = errors.New("node not found")
	ErrPairingTokenInvalid   = errors.New("pairing token invalid")
	ErrPairingTokenExpired   = errors.New("pairing token expired")
	ErrAgentTokenInvalid     = errors.New("agent token invalid")
	ErrNodeAlreadyRegistered = errors.New("node already registered")
)
