package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) Login(ctx echo.Context) error {
    return ctx.JSON(http.StatusNotImplemented, echo.Map{"message": "not implemented"})
}

func (s *Server) Signup(ctx echo.Context) error {
    return ctx.JSON(http.StatusNotImplemented, echo.Map{"message": "not implemented"})
}