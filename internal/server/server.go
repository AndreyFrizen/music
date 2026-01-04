package server

import "mess/internal/repository/store"

type Server struct {
	Store *store.Store
}

func NewServer(store *store.Store) *Server {
	return &Server{
		Store: store,
	}
}
