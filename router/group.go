package router

import (
	"net/http"

	"uprav/api"
	"uprav/converter"

	"github.com/labstack/echo/v4"

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

	response, err := converter.ConvertGroup[[]api.GroupsResponse](groups)
  if err != nil {
    return e.JSON(http.StatusInternalServerError, api.InternalServerError{Message: ptrString("failed to build response")})
  }

  return e.JSON(http.StatusOK, response)
}