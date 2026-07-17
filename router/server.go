package router

import (
	"context"

	"uprav/repository/gormrepository"

)

// Server 構造体：ここにDB接続やサービスをまとめます
type Server struct {
	ctx  context.Context
	taskRepo gormrepository.TaskRepository // インターフェースで持つ
	userRepo gormrepository.UserRepository // インターフェースで持つ
	groupRepo gormrepository.GroupRepository // インターフェースで持つ
	// repo, service などが必要ならここに追加していきます
}

// NewServer：main.goから呼ばれる初期化関数
func NewServer(
		ctx context.Context, 
		taskRepo gormrepository.TaskRepository, 
		userRepo gormrepository.UserRepository,
		groupRepo gormrepository.GroupRepository,
	) *Server {

	return &Server{
		ctx:  ctx,
		taskRepo: taskRepo,
		userRepo: userRepo,
		groupRepo: groupRepo,
	}
}
