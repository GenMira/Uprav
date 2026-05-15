package router

import (
	"net/http"

	"uprav/api"
	"uprav/converter"
	"uprav/model"

	"github.com/labstack/echo/v4"
)

// --- 以下、openapi.yaml で定義した operationId をメソッドとして実装 ---

// GetTasks (GET /api/tasks) の実装
func (s *Server) GetTasks(e echo.Context) error {
	tasks, err := s.taskRepo.GetAllTasks(e.Request().Context())
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get tasks")})
	}

	response, err := converter.Convert[[]api.TaskResponse](tasks)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to convert tasks")})
	}

	return e.JSON(http.StatusOK, response)
}

// CreateTask (POST /api/newtask) の実装
func (s *Server) CreateTask(e echo.Context) error {
	var req api.TaskRequest

	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("invalid request body")})
	}

	if req.Name == "" {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("name is required")})
	}

	task, err := converter.Convert[model.Task](req)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to convert request")})
	}

	if err := s.taskRepo.CreateTask(e.Request().Context(), &task); err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to create task")})
	}

	response, err := converter.Convert[api.TaskResponse](task)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
	}

	return e.JSON(http.StatusCreated, response)
}

func ptrString(s string) *string {
	return &s
}