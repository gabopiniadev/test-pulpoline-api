package errors

import "errors"

var (
	ErrQueueClosed   = errors.New("la cola está cerrada")
	ErrQueueFull     = errors.New("la cola está llena")
	ErrMissingAPIKey = errors.New("API key no está configurada")
)
