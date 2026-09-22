package organizations

import "errors"

var (
	ErrMemberNotFound = errors.New("organization member not found")
	ErrForbidden      = errors.New("not authorized for this organization")
	ErrInvalidRole    = errors.New("invalid role")
)
