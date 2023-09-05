package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Server is a struct for connection handlers and database
type Server struct {
	listenAddr string
}

// NewServer is a default constructor
func NewServer(listenAdder string) *Server {
	return &Server{
		listenAddr: listenAdder,
	}
}

// Run configures router and starts Server
func (s *Server) Run() error {
	e := echo.New()

	e.GET("/tokens", generateToken)
	e.PATCH("/tokens/refresh", updateToken)

	return http.ListenAndServe(s.listenAddr, e)
}

// generateToken gets user Id and returns a pair of tokens
func generateToken(c echo.Context) error {
	return nil
}

// updateToken updates Access token
func updateToken(c echo.Context) error {
	return nil
}
