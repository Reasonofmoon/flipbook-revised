package models

import "time"

const (
	StatusPending      = "pending"
	StatusConverting   = "converting"
	StatusReady        = "ready"
	StatusError        = "error"
	StatusRegenerating = "regenerating"
)

type Flipbook struct {
	ID           string
	Title        string
	Slug         string
	Description  string
	Client       string // academy the flipbook was made for (agency model)
	Filename     string
	FileSize     int64
	PageCount    int
	Status       string
	ErrorMessage string
	PageWidth    int
	PageHeight   int
	IsPublic     bool   // false: only reachable with AccessKey in the link
	AccessKey    string // random key required as ?k= when not public
	GridFSFileID string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	ConvertedAt  *time.Time
}
