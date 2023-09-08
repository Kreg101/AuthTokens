package api

import (
	"encoding/base64"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Repository interface represents a storage for
// saving Refresh tokens
type Repository interface {
	InsertRefresh(refreshToken string) error
	CheckRefresh(oldRefreshToken, newRefreshToken string) (bool, error)
}

// Server is a struct for connection handlers and database
type Server struct {
	listenAddr string
	storage    Repository
}

// NewServer is a default constructor
func NewServer(listenAdder string, storage Repository) *Server {
	return &Server{
		listenAddr: listenAdder,
		storage:    storage,
	}
}

// Run configures router and starts Server
func (s *Server) Run() error {
	e := echo.New()

	e.GET("/tokens", s.getTokens)
	e.POST("/tokens/refresh", s.updateTokens)

	return http.ListenAndServe(s.listenAddr, e)
}

// generateToken is a handler, that gets GUID and
// returns a pair of Access and Refresh tokens
func (s *Server) getTokens(c echo.Context) error {
	guid := c.QueryParam("guid")
	if guid == "" {
		return c.JSON(http.StatusBadRequest,
			Msg{"invalid empty GUID"})
	}

	tokens, err := generateTokenPair(guid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			Msg{"can't generate tokens"})
	}

	err = s.storage.InsertRefresh(tokens["refresh"])
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			Msg{"can't save refresh token"})
	}

	tokens["refresh"] = base64.StdEncoding.EncodeToString([]byte(tokens["refresh"]))

	return c.JSON(http.StatusOK, tokens)
}

// updateToken is a handler, that updates a pair of Access and Refresh tokens
func (s *Server) updateTokens(c echo.Context) error {
	var refresh Token
	err := c.Bind(&refresh)
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			Msg{"can't unmarshal json"})
	}

	tokens, err := generateTokenPair(refresh.Refresh)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			Msg{"can't generate tokens"})
	}

	exists, err := s.storage.CheckRefresh(refresh.Refresh, tokens["refresh"])
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			Msg{"can't check tokens existence"})
	}

	if !exists {
		return c.JSON(http.StatusNotFound,
			Msg{"token is invalid"})
	}

	tokens["refresh"] = base64.StdEncoding.EncodeToString([]byte(tokens["refresh"]))

	return c.JSON(http.StatusOK, tokens)
}

// Msg represents message to user
type Msg struct {
	Err string `json:"error"`
}

// Token for binding
type Token struct {
	Refresh string `json:"refresh"`
}
