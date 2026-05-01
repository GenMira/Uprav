package router

import (
	"context"

	"uprav/repository/gormrepository"

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
