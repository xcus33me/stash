package file

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID            uuid.UUID
	RoomID        uuid.UUID
	OriginalName  string
	UploaderAlias *string
	SizeBytes     int64
	MimeType      string
	StorageKey    string
	SHA256        string
	ExpiresAt     *time.Time
	UploadedAt    time.Time
	DeletedAt     *time.Time
}

func (f *File) IsDeleted() bool {
	return f.DeletedAt != nil
}

func (f *File) IsExpired() bool {
	return f.ExpiresAt != nil && time.Now().After(*f.ExpiresAt)
}
