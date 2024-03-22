package core

import "errors"

var (
	ErrResponseNotImplemented   = errors.New("response type not implemented")
	ErrInvalidRouteResponseType = errors.New("invalid route response type")
)
