package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"

	"net/http"
	"uprav/api"
	"uprav/repository/gormrepository"
	"uprav/router"
	"uprav/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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
	if user == "" || password == "" || host == "" || port == "" || database == "" {
		log.Fatal("missing required database environment variables: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME")
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
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.POST, echo.OPTIONS},
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

	repo, err := gormrepository.NewGormRepository(db)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	api.RegisterHandlers(
		e,
		router.NewServer(ctx, repo),
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