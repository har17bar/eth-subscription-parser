package rest

import (
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(s *Subscription, t *Transactions) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/subscribe", s.Subscribe)
	mux.HandleFunc("/transactions", t.GetTransactions)
	mux.HandleFunc("/current", t.GetCurrentProcessedBlock)

	return &Server{
		mux: mux,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
