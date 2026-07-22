package router

import (
	"database/sql"
	"errors"
	"net/http"

	"uprav/api"
	"uprav/converter"
	"uprav/model"

	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"


)

func (s *Server) GetGroups(e echo.Context) error {
	loginUID, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	groups,err := s.groupRepo.GetGroups(e.Request().Context(),loginUID)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get groups")})
	}

	if groups == nil {
		groups = []model.Group{}
	}

	response, err := converter.ConvertGroup[[]api.GroupResponse](groups)
  if err != nil {
    return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
  }

  return e.JSON(http.StatusOK, response)
}

func (s *Server) CreateGroup(e echo.Context) error {
	_, _ , err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	var req api.GroupRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("invalid request body")})
	}

	group, err := s.groupRepo.CreateGroup(e.Request().Context(),req.Name,req.Members)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to create group")})
	}

	response, err := converter.ConvertGroup[api.GroupResponse](group)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
	}
	return e.JSON(http.StatusCreated, response)
}

func (s *Server) UpdateGroup(e echo.Context,id openapi_types.UUID) error {
	_, _, err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	groupID := uuid.UUID(id)
	var req api.UpdateGroupRequest
	if err := e.Bind(&req); err != nil {
		return e.JSON(http.StatusBadRequest, api.BadRequest{Message: ptrString("invalid request body")})
	}

	updatedGroup, err := s.groupRepo.UpdateGroup(e.Request().Context(), groupID, req.Name, req.Members)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.JSON(http.StatusNotFound, api.NotFound{Message: ptrString("group not found")})
		}
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to update group")})
	}

	response, err := converter.ConvertGroup[api.GroupResponse](updatedGroup)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
	}

	return e.JSON(http.StatusOK, response)
}

func (s *Server) DeleteGroup(e echo.Context,id openapi_types.UUID) error {
	_, _, err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	groupID := uuid.UUID(id)
	err = s.groupRepo.DeleteGroup(e.Request().Context(), groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.JSON(http.StatusNotFound, api.NotFound{Message: ptrString("group not found")})
		}
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to delete group")})
	}

	return e.NoContent(http.StatusNoContent)
}

func (s *Server) GetGroup(e echo.Context,	id openapi_types.UUID) error {
	_, _, err := GetDataFromToken(e)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, api.BadRequest{Message: ptrString(err.Error())})
	}

	groupID := uuid.UUID(id)
	group, err := s.groupRepo.GetGroup(e.Request().Context(), groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.JSON(http.StatusNotFound, api.NotFound{Message: ptrString("group not found")})
		}
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to get group")})
	}

	response, err := converter.ConvertGroup[api.GroupResponse](group)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
	}

	return e.JSON(http.StatusOK, response)
}
