package router

import (
	"net/http"
	"uprav/api"
	"uprav/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) Login(e echo.Context) error {
	var req api.AuthenticationRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("Invalid request format")})
	}

	if req.Name == "" || req.Password == "" {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("Missing required fields")})
	}

	user, err := s.userRepo.GetUserByName(e.Request().Context(), req.Name)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString("Invalid name or password")})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString("Invalid name or password")})
	}

	token := "dummy-token-for-" + user.Name

	return e.JSON(http.StatusOK, api.AuthenticationResponse{
		Uid:   &user.UID,
		Name:  &user.Name,
		Token: &token,
	})
}

func (s *Server) Signup(e echo.Context) error {
	var req api.AuthenticationRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("Invalid request format")})
	}

	if req.Name == "" || req.Password == "" {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("Missing required fields")})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("Failed to process password")})
	}

	user := model.User{
		Name:     req.Name,
		Password: string(hashedPassword),
	}
	if err := s.userRepo.CreateUser(e.Request().Context(), &user); err != nil {
		// すでに同じUIDが存在する場合などのエラーハンドリング
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("Could not create user")})
	}

	// ダミートークンの生成
	token := "dummy-token-for-" + user.Name

	return e.JSON(http.StatusOK, api.AuthenticationResponse{
		Uid:  &user.UID,
		Name: &user.Name,
		Token: &token,
	})
}