package converter

import (
	"errors"

	"uprav/api"
	"uprav/model"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Convert converts known types between API and model representations.
func Convert[T any](src any) (T, error) {
	var zero T
	switch v := src.(type) {
	case []model.Task:
		switch any(zero).(type) {
		case []api.TaskResponse:
			out := make([]api.TaskResponse, len(v))
			for i, task := range v {
				out[i] = convertTaskToResponse(task)
			}
			return any(out).(T), nil
		}
	case model.Task:
		switch any(zero).(type) {
		case api.TaskResponse:
			return any(convertTaskToResponse(v)).(T), nil
		}
	case api.NewTaskRequest:
		switch any(zero).(type) {
		case model.Task:
			return any(convertRequestToModel(v)).(T), nil
		}
	case *api.NewTaskRequest:
		switch any(zero).(type) {
		case model.Task:
			return any(convertRequestToModel(*v)).(T), nil
		}
	case api.UpdateTaskRequest:
		switch any(zero).(type) {
		case model.Task:
			return any(convertUpdateRequestToModel(v)).(T), nil
		}
	}

	return zero, errors.New("unsupported conversion")
}

func convertTaskToResponse(task model.Task) api.TaskResponse {
	resp := api.TaskResponse{}

	resp.Name = task.Name
	resp.Priority = task.Priority
	resp.Deadline = task.Deadline
	resp.IsEveryday = task.IsEveryday
	if !task.Period.IsZero() {
		resp.Period = &task.Period
	}
	if task.Assign != uuid.Nil {
		assign := openapi_types.UUID(task.Assign)
		resp.Assign = &assign
	}
	if task.Group != uuid.Nil {
		group := openapi_types.UUID(task.Group)
		resp.Group = &group
	}
	if task.Tag != "" {
		resp.Tag = &task.Tag
	}
	if task.Description != "" {
		resp.Description = &task.Description
	}

	// if len(task.Assumption) > 0 {
	// 	assumption := make([]openapi_types.UUID, len(task.Assumption))
	// 	for i, id := range task.Assumption {
	// 		assumption[i] = openapi_types.UUID(id)
	// 	}
	// 	resp.Assumption = &assumption
	// }

	return resp
}

func convertRequestToModel(req api.NewTaskRequest) model.Task {
	task := model.Task{}

	task.Name = req.Name
	task.Priority = req.Priority
	task.Deadline = req.Deadline
	task.IsEveryday = req.IsEveryday
	if req.Period != nil {
		task.Period = *req.Period
	}
	if req.Assign != nil {
		task.Assign = *req.Assign
	}
	if req.Group != nil {
		task.Group = *req.Group
	}
	if req.Tag != nil {
		task.Tag = *req.Tag
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	// if req.Assumption != nil {
	// 	assumption := make([]uuid.UUID, len(*req.Assumption))
	// 	for i, id := range *req.Assumption {
	// 		assumption[i] = uuid.UUID(id)
	// 	}
	// 	task.Assumption = assumption
	// }

	return task
}

func convertUpdateRequestToModel(req api.UpdateTaskRequest) model.Task {
	task := model.Task{}

	task.Name = *req.Name
	task.Priority = *req.Priority
	task.Tag = *req.Tag
	task.Deadline = *req.Deadline
	task.Period = *req.Period
	task.IsEveryday = *req.IsEveryday
	task.Assign = *req.Assign
	task.Group = *req.Group
	// if req.Assumption != nil {
	// 	assumption := make([]uuid.UUID, len(*req.Assumption))
	// 	for i, id := range *req.Assumption {
	// 		assumption[i] = uuid.UUID(id)
	// 	}
	// 	task.Assumption = assumption
	// }

	return task
}