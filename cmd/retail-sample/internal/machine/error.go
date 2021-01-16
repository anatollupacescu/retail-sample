package usecase

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
)
