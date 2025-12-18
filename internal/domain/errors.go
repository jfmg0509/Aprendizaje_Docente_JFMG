package domain

import "errors"

var (
    ErrNotFound        = errors.New("not found")
    ErrDuplicate       = errors.New("duplicate")
    ErrValidation      = errors.New("validation error")
    ErrInvalidRole     = errors.New("invalid role")
    ErrInvalidAccess   = errors.New("invalid access type")
    ErrInactiveEntity  = errors.New("entity is inactive")
)
