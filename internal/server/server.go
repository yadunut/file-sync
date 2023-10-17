package server

import (
	"fmt"
	"path/filepath"
	"sync"

	"slices"

	"github.com/yadunut/file-sync/internal/contracts"
	"github.com/yadunut/file-sync/internal/server/db"
	"github.com/yadunut/file-sync/internal/util"
	"go.uber.org/zap"
)

type Server struct {
	Config util.Config
	Db     *db.DB
	Log    *zap.SugaredLogger
}

func CreateServer(Db *db.DB, log *zap.SugaredLogger, config util.Config) *Server {
	return &Server{
		Config: config,
		Db:     Db,
		Log:    log,
	}
}

func (s *Server) Start(wg *sync.WaitGroup) {
	// start the background processing routines
}

func (s *Server) Version() contracts.VersionRes {
	return contracts.VersionRes{Version: util.VERSION, Success: true}
}

func (s *Server) WatchUp(req contracts.WatchUpReq) (contracts.WatchUpRes, error) {
	if !filepath.IsAbs(req.Path) {
		return contracts.WatchUpRes{Success: false}, fmt.Errorf("Path must be absolute")
	}

	dirs, err := s.Db.GetDirectories()
	if err != nil {
		return contracts.WatchUpRes{Success: false}, err
	}
	if dirs != nil {
		if slices.ContainsFunc(dirs, func(entry db.Directory) bool { return util.IsDirChild(req.Path, entry.Path) }) {
			return contracts.WatchUpRes{Success: false}, fmt.Errorf("Path is already being watched")
		}
	}
	s.Db.AddDirectory(req.Path)
	return contracts.WatchUpRes{Success: true}, nil
}

func (s *Server) WatchDown(req contracts.WatchDownReq) (contracts.WatchDownRes, error) {
	if !filepath.IsAbs(req.Path) {
		return contracts.WatchDownRes{Success: false}, fmt.Errorf("Path must be absolute")
	}
	err := s.Db.RemoveDirectory(req.Path)
	if err != nil {
		return contracts.WatchDownRes{Success: false}, fmt.Errorf("Path must be absolute")
	}
	return contracts.WatchDownRes{Success: true}, nil
}

func (s *Server) WatchList() (contracts.WatchListRes, error) {
	dirs, err := s.Db.GetDirectories()
	if err != nil {
		return contracts.WatchListRes{Success: false}, err
	}
	return contracts.WatchListRes{Success: true, Directories: dirs}, nil
}
