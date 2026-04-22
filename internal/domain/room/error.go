package room

import "errors"

var (
	ErrNotFound      = errors.New("room not found")
	ErrLocked        = errors.New("room is locked")
	ErrInvalidOwner  = errors.New("invalid owners token")
	ErrTitleRequired = errors.New("room title required")
)
