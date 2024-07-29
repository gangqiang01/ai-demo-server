package constants

import (
	"errors"
)

var (
	ErrInvalidParms          = errors.New("Invalid parms")
	ErrInvalidOperation      = errors.New("Invalid operation!")
	ErrInvalidRequest        = errors.New("Invalid or nil request!")
	ErrChannelClosed         = errors.New("channel has been closed!")
	ErrEmptyContent          = errors.New("empty contents!")
	ErrInvalidResponseStruct = errors.New("invalid response structure.")
	ErrNoSuchDevice          = errors.New("No Such Device.")
	ErrUnknown               = errors.New("unknown error.")
	ErrCreateDeviceFailed    = errors.New("create device failed.")
	ErrDeviceNotExists       = errors.New("device is not exists.")
	ErrNotSupport            = errors.New("Not support.")
	ErrInvalidURLFormat      = errors.New("invalid url format.")
	ErrHasExists             = errors.New("item has already exists.")
	ErrStackNotReady         = errors.New("stack not ready")
)
