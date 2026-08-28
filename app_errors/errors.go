package apperrors

import "errors"

var ErrNotFound error = errors.New("quote not found")
var ErrAPINotAvailable error = errors.New("failed to connect thirg-party service")
var ErrInvalidPair error = errors.New("invalid pair. Choose another one")
var ErrInvalidPairFormat = errors.New("invalid currency pair format. Use format USDGBP")
