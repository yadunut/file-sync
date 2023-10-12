package server

import (
	"log"
	"net/http"

	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yadunut/file-sync/internal/contracts"
	"github.com/yadunut/file-sync/internal/server/db"
	"github.com/yadunut/file-sync/internal/util"
)

type Server struct {
	config     contracts.Config
	HttpServer *http.Server
	router     chi.Router
	Db         *db.DB
	Log        *log.Logger
}

type Routes map[string]http.HandlerFunc

func CreateServer(Db *db.DB, log *log.Logger) *Server {
	router := chi.NewRouter()
	config := contracts.Config{Port: ":8080", Host: "localhost"}
	server := &http.Server{Addr: config.GetUrl(), Handler: router}
	return &Server{
		config:     config,
		HttpServer: server,
		router:     router,
		Db:         Db,
		Log:        log,
	}
}

func (s *Server) Start() error {
	s.router.Use(middleware.Logger)
	s.router.Get("/version", s.VersionFunc)
	return s.HttpServer.ListenAndServe()
}

func (s *Server) VersionFunc(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(contracts.Version{Version: util.VERSION})
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
