package repository

import "errors"

var (
	ErrDuplicateDetails = errors.New("username or email already exists")
)
