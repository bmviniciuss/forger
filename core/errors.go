package core

import "errors"

var (
	ErrResponseNotImplemented         = errors.New("response type not implemented")
	ErrInvalidRouteResponseType       = errors.New("invalid route response type")
	ErrInvalidRequestBodyKeysAmount   = errors.New("requestBody function only supports one or zero aguments")
	ErrInvalidRequestBodyArgumentType = errors.New("requestBody argument should be a string")
)
