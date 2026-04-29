package router

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"uprav/repository/gormrepository"
	//"Uprav/api" // go.modのmodule名が "uprav" の場合。適宜書き換えてください
)

// Server 構造体：ここにDB接続やサービスをまとめます
type Server struct {
	ctx  context.Context
	repo *gormrepository.Repository
	// repo, service などが必要ならここに追加していきます
}

// NewServer：main.goから呼ばれる初期化関数
func NewServer(ctx context.Context, repo *gormrepository.Repository) *Server {
	return &Server{
		ctx:  ctx,
		repo: repo,
	}
}

// --- 以下、openapi.yaml で定義した operationId をメソッドとして実装 ---
// これを 1 つでも書き忘れると、api.RegisterHandlers でエラーになります

// GetTasks (GET /api/tasks) の実装
func (s *Server) GetTasks(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, []string{"task1", "task2"})
}

// CreateTask (POST /api/newtask) の実装
func (s *Server) CreateTask(ctx echo.Context) error {
	return ctx.JSON(http.StatusCreated, map[string]string{"message": "created!"})
}