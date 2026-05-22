package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"net/http"
	"uprav/api"
	"uprav/repository/gormrepository"
	"uprav/router"
	"uprav/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main(){

	e := echo.New()
	user := os.Getenv("NS_MARIADB_USER")
	password := os.Getenv("NS_MARIADB_PASSWORD")
	host := os.Getenv("NS_MARIADB_HOSTNAME")
	port := os.Getenv("NS_MARIADB_PORT")
	database := os.Getenv("NS_MARIADB_DATABASE")
	jwtSecretEnv := []byte(os.Getenv("JWT_SECRET_KEY"))
	if user == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatal("missing required database environment variables: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME")
	}
	if jwtSecretEnv == nil {
		log.Fatal("missing required environment variable: JWT_SECRET_KEY")
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FTokyo",
		user,
		password,
		host,
		port,
		database,
	)
	gormLogLevel := logger.Silent

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:8080",
		},
		AllowMethods: []string{
      echo.GET, 
      echo.POST, 
      echo.PUT, 
      echo.PATCH, 
      echo.DELETE, 
      echo.OPTIONS,
    },
    AllowHeaders: []string{
      echo.HeaderOrigin,
      echo.HeaderContentType,
      echo.HeaderAccept,
      echo.HeaderAuthorization,
    },
    AllowCredentials: true,
  }))

	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecretEnv),
		Skipper: func(c echo.Context) bool {
			// CORSのPreflightリクエスト(OPTIONS)はJWTチェックをスキップ
			if c.Request().Method == echo.OPTIONS {
				return true
			}
			// パスのブレをなくすため、リクエストの生パスで判定
			path := c.Request().URL.Path
			if path == "/api/login" || path == "/api/signup" || strings.HasSuffix(path, "/login") || strings.HasSuffix(path, "/signup") {
				return true
			}
			return false
		},
		ErrorHandler: func(c echo.Context, err error) error {
			fmt.Printf("[JWT ERROR] Authentication failed: %v\n", err)
			return c.JSON(
				http.StatusUnauthorized,
				map[string]string{
					"message": fmt.Sprintf("Auth failed: %v", err.Error()),
				})
		},
	}))


	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(gormLogLevel),
		TranslateError: true,
	})

	if err := model.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	if err != nil {
		log.Fatal(err)
	}

	taskRepo := gormrepository.NewTaskRepository(db)
  userRepo := gormrepository.NewUserRepository(db)

	ctx := context.Background()

	api.RegisterHandlers(
		e,
		router.NewServer(
			ctx,
			taskRepo,
			userRepo,
		),
	)
	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: e,
	}

	go func() {
		if err := e.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// シグナル待ちやクリーンアップ処理を追加
	select {}

}