package kolosal

import "errors"

var (
	errNoMessageProvided      = errors.New("no message provided")
	errEmptyResponseFromModel = errors.New("empty response from model")
)
