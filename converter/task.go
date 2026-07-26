package converter

import (
	"errors"

	"uprav/api"
	"uprav/model"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Convert converts known types between API and model representations.
// 更新系（UpdateTaskRequest）のようにURLパスのIDが必要な場合、第2引数に uint 型の ID を渡すことができます。
func Convert[T any](src any, id ...uint) (T, error) {
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
	case *model.Task:
		switch any(zero).(type) {
		case api.TaskResponse:
			// v はポインタなので、*v でデリファレンスして実体を渡す
			return any(convertTaskToResponse(*v)).(T), nil
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
			// 安全にパスから渡ってきた ID を利用する（未指定の場合は 0）
			var targetID uint
			if len(id) > 0 {
				targetID = id[0]
			}
			return any(convertUpdateRequestToModel(targetID, v)).(T), nil
		}
	}

	return zero, errors.New("unsupported conversion")
}

func convertTaskToResponse(task model.Task) api.TaskResponse {
	resp := api.TaskResponse{}

	resp.Id = ptrInt(int(task.ID))
	resp.Name = task.Name
	resp.Priority = task.Priority
	resp.Deadline = task.Deadline
	resp.IsEveryday = task.IsEveryday

	if !task.Period.IsZero() {
		resp.Period = &task.Period
	}
	if task.Assigner != "" {
		resp.Assigner = &task.Assigner
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

	return resp
}

//新規作成用
func convertRequestToModel(req api.NewTaskRequest) model.Task {
	task := model.Task{}

	task.Name = req.Name
	task.Priority = req.Priority
	task.Deadline = req.Deadline
	task.IsEveryday = req.IsEveryday

	if req.Period != nil {
		task.Period = *req.Period
	}
	if req.Assignee != nil {
    task.AssignGroup = *req.AssignGroup
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

	return task
}

func convertUpdateRequestToModel(id uint, req api.UpdateTaskRequest) model.Task {
	task := model.Task{}

	task.ID = id

	if req.Name != nil {
		task.Name = *req.Name
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.Deadline != nil {
		task.Deadline = *req.Deadline
	}
	if req.IsEveryday != nil {
		task.IsEveryday = *req.IsEveryday
	}
	if req.Period != nil {
		task.Period = *req.Period
	}
	if req.Tag != nil {
		task.Tag = *req.Tag
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Assignee != nil {
		task.AssignGroup = *req.AssignGroup
	}
	if req.Group != nil {
		task.Group = *req.Group
	}

	return task
}

func ptrInt(i int) *int {
  return &i
}