package repository

import "errors"

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")
var ErrInvalidRef = errors.New("invalid_reference")
var ErrValidation = errors.New("validation_error")
