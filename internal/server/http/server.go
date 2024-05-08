package http

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yadunut/file-sync/internal/contracts"
	"github.com/yadunut/file-sync/internal/server"
)

type HttpServer struct {
	Server     *server.Server
	router     chi.Router
	httpServer http.Server
}

func NewHttpServer(s *server.Server) *HttpServer {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	s.Log.Info("Starting server on ", s.Config.GetUrl())
	return &HttpServer{
		Server:     s,
		router:     router,
		httpServer: http.Server{Addr: s.Config.GetUrl(), Handler: router},
	}
}

func (s *HttpServer) Start() error {
	var wg sync.WaitGroup
	s.router.Get("/version", s.VersionFunc)
	s.router.Post("/watch/up", s.WatchUpFunc)
	s.router.Post("/watch/down", s.WatchDownFunc)
	s.router.Get("/watch", s.WatchListFunc)
	s.Server.Start(&wg)
	wg.Add(1)
	go func() { s.httpServer.ListenAndServe(); wg.Done() }()
	wg.Wait()
	return nil
}

func (s *HttpServer) VersionFunc(w http.ResponseWriter, r *http.Request) {
	writeSuccess(w, s.Server.Version())
}

func (s *HttpServer) WatchUpFunc(w http.ResponseWriter, r *http.Request) {
	var watchUp contracts.WatchUpReq
	err := json.NewDecoder(r.Body).Decode(&watchUp)
	if err != nil {
		writeError(w, err)
		return
	}
	res, err := s.Server.WatchUp(watchUp)
	if err != nil {
		writeError(w, err)
		return
	}
	if res.Success {
		writeSuccess(w, res)
	} else {
		writeFailure(w, res)
	}
}

func (s *HttpServer) WatchDownFunc(w http.ResponseWriter, r *http.Request) {
	var watchDown contracts.WatchDownReq
	err := json.NewDecoder(r.Body).Decode(&watchDown)
	if err != nil {
		s.Log.Error("error decoding json", err)
		writeError(w, err)
	}
	res, err := s.WatchDown(watchDown)
	if err != nil {
		writeError(w, err)
		return
	}

	if res.Success {
		writeSuccess(w, res)
	} else {
		writeFailure(w, res)
	}
}

func (s *HttpServer) WatchListFunc(w http.ResponseWriter, r *http.Request) {
	res, err := s.WatchList()
	if err != nil {
		writeError(w, err)
		return
	}
	if res.Success {
		writeSuccess(w, res)
	} else {
		writeFailure(w, res)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(contracts.ErrorRes{Error: err.Error(), Success: false})
}

func writeSuccess(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(v)
}

func writeFailure(w http.ResponseWriter, v any) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(v)
}
