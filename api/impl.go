package api

import "net/http"

type Server struct{}

func (s Server) GetCoveragePath(w http.ResponseWriter, r *http.Request, path string, params GetCoveragePathParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PostCoveragePath(w http.ResponseWriter, r *http.Request, path string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
