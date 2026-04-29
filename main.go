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

	"github.com/labstack/echo/v4"

	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main(){

	e := echo.New()
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FTokyo",
		user,
		password,
		host,
		port,
		database,
	)
	gormLogLevel := logger.Silent


	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(gormLogLevel),
		TranslateError: true,
	})

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