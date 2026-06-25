package router

import (
	"net/http"
	"os"
	"time"
	"uprav/api"
	"uprav/model"

	"github.com/go-task/task/v3/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)


type JwtCustomClaims struct {
	UID  int    `json:"uid"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func generateToken(uid int, name string) (string, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

	claims := &JwtCustomClaims{
		UID:  uid,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			//有効時間は24時間
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 署名アルゴリズム（HS256）を指定してトークンオブジェクトを作成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 秘密鍵で署名して文字列のトークンを生成
	return token.SignedString(jwtSecret)
}

func GetDataFromToken(e echo.Context) (int, string, error) {
	token, ok := e.Get("user").(*jwt.Token)
	if !ok {
		return 0, "", errors.New("Unauthorized: Token missing")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("Unauthorized: Invalid claims")
	}

	uidFloat, ok := claims["uid"].(float64)
	if !ok {
		return 0, "", errors.New("Unauthorized: UID missing in token")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return 0, "", errors.New("Unauthorized: Name missing in token")
	}

	return int(uidFloat), name, nil
}

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

	token, err := generateToken(user.UID, user.Name)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("Failed to generate token")})
	}

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
	token, err := generateToken(user.UID, user.Name)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("Failed to generate token")})
	}

	return e.JSON(http.StatusOK, api.AuthenticationResponse{
		Uid:  &user.UID,
		Name: &user.Name,
		Token: &token,
	})
}

func (s *Server) GetCurrentUser(e echo.Context) error {
	loginUID, loginName, err := GetDataFromToken(e)
	if err != nil {
		// トークンが無効、または期限切れの場合
		return e.JSON(http.StatusUnauthorized, api.BadRequest{
			Message: ptrString("token expired or invalid"),
		})
	}

	// トークンが有効であれば、現在のユーザー情報をそのまま返す
	response := struct {
		UID  int    `json:"uid"`
		Name string `json:"name"`
	}{
		UID:  loginUID,
		Name: loginName,
	}

	return e.JSON(http.StatusOK, response)
}