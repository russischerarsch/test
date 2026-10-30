package domain

import "errors"

var ErrNotFound = errors.New("quote not found")
var ErrAPINotAvailable = errors.New("failed to connect thirg-party service")
var ErrInvalidPairFormat = errors.New("invalid currency pair format. Use format USDGBP")
var ErrAlreadyExists = errors.New("pair already exists")
