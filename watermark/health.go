package main

import (
    "net/http"
    "log"
)

func (s *Server) handleHealth() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        respond(w, r, http.StatusOK, nil)
    }
}
