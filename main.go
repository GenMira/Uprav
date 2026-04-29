package main

import (
    "uprav/api"
		"errors"
    "uprav/router"
		"uprav/repository/gormrepository"
    "net/http"
    "github.com/labstack/echo/v4"

		"gorm.io/gorm"
    "gorm.io/driver/mysql"


)

func main() {
	e := echo.New()

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

	api.RegisterHandlers(
		e,
		router.NewServer(ctx, repo, activityService, notificationService, traqService, isDev),
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

	e.Logger.Fatal(e.Start(":1124"))
	//localhost:1124
}