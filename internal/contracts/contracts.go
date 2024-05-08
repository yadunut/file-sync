package contracts

import (
	"time"

	"github.com/yadunut/file-sync/internal/server/db"
)

type ErrorRes struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type VersionRes struct {
	Version string `json:"version"`
	Success bool   `json:"success"`
}

// WatchUpReq expects the full path. if its not a full path, it will error
type WatchUpReq struct {
	Path string `json:"path"`
}

type WatchUpRes struct {
	Success bool `json:"success"`
}

type WatchDownReq struct {
	Path string `json:"path"`
}

type WatchDownRes struct {
	Success bool `json:"success"`
}

type WatchListRes struct {
	Directories []db.Directory `json:"directories"`
	Success     bool           `json:"success"`
}
