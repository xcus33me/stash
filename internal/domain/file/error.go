package file

import "errors"

var (
	ErrNotFound       = errors.New("file not found")
	ErrExpired        = errors.New("file expired")
	ErrDeleted        = errors.New("file deleted")
	ErrNameRequired   = errors.New("file name required")
	ErrInvalidSize    = errors.New("invalid file size")
	ErrQuotaExceeded  = errors.New("room storage quota exceeded")
	ErrChecksumFailed = errors.New("checksum mismatch")
)
