package router

import (
	"net/http"

	"errors"
	"fmt"
	"uprav/api"
	"uprav/converter"
	"uprav/model"
	"log"

	"gorm.io/gorm"

	//"github.com/docker/docker/libnetwork/drivers/null"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ptrString(s string) *string {
	return &s
}

func (s *Server) GetTasks(e echo.Context) error {
	loginUID, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	tasks, err := s.taskRepo.GetAllTasks(e.Request().Context(),loginUID)
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
	//tokenの確認
	loginUID, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	var req api.NewTaskRequest

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

	if task.AssignGroup!=uuid.Nil{
		task.Uid = *req.Assignee
		user,err := s.userRepo.GetUserByUID(e.Request().Context(), task.Uid)
		if err != nil{
			return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get username from UID")})
		}
		task.Assigner = fmt.Sprintf("%s(%d)", user.Name, task.Uid)
	}else{
		task.Uid = loginUID
	}

	if err := s.taskRepo.CreateTask(e.Request().Context(), &task); err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to create task")})
	}

	response, err := converter.Convert[api.TaskResponse](task)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
	}

	if task.AssignGroup != uuid.Nil {
		group, err := s.groupRepo.GetGroup(e.Request().Context(), task.AssignGroup)
		if err != nil {
			return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get group from ID")})
		}
		response.AssignGroup = &group.Name
		log.Printf("[CreateTask] response.AssignGroup: %d,task.AssignGroup: %s",response.AssignGroup,task.AssignGroup)
	}

	return e.JSON(http.StatusCreated, response)
}

// 引数に `id int` が自動で追加される（oapi-codegenの仕様）
func (s *Server) DeleteTask(e echo.Context, id int) error {
	loginUID, _, err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}
	if err := s.taskRepo.DeleteTask(e.Request().Context(), id, loginUID); err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to delete task")})
	}

	return e.JSON(http.StatusOK, api.Accepted{Message: ptrString("Task was deleted successfully")})
}

func (s *Server) UpdateTask(e echo.Context, id int) error {
	loginUID, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	var req api.UpdateTaskRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("Invalid format")})
	}

	task, err := converter.Convert[model.Task](req)
	task.ID = uint(id)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("Failed to convert request")})
	}
	task.Uid = loginUID
	
	if err := s.taskRepo.UpdateTask(e.Request().Context(), &task); err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("Failed to update task")})
	}

	response, err := converter.Convert[api.TaskResponse](task)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
	}

	return e.JSON(http.StatusCreated, response)
}

func (s *Server) GetTask(e echo.Context, id int) error {
	loginUID, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	task, err := s.taskRepo.GetTask(e.Request().Context(), id, loginUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.JSON(http.StatusNotFound, api.NotFound{Message: ptrString("task not found")})
		}
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get task")})
	}

	response, err := converter.Convert[api.TaskResponse](task)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString(err.Error())})
	}

	return e.JSON(http.StatusOK, response)
}

func (s *Server) GetTags(e echo.Context) error {
	loginUID, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	tags, err := s.taskRepo.GetTags(e.Request().Context(), loginUID)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get tags")})
	}

	response := api.TagsResponse{
		Tags: &tags,
	}

	return e.JSON(http.StatusOK, response)
}

