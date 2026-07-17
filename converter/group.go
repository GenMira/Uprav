package converter

import (
	"errors"
	"uprav/api"
	"uprav/model"
)

// ConvertGroup は Group 系の型マッピングを行う型スイッチです（既存の Convert 関数内にマージしてもOKです）
func ConvertGroup[T any](src any) (T, error) {
	var zero T
	switch v := src.(type) {
	case []model.Group:
		switch any(zero).(type) {
		case []api.GroupsResponse:
			out := make([]api.GroupsResponse, len(v))
			for i, g := range v {
				out[i] = convertGroupToResponse(g)
			}
			return any(out).(T), nil
		}
	case model.Group:
		switch any(zero).(type) {
		case api.GroupsResponse:
			return any(convertGroupToResponse(v)).(T), nil
		}
	case *model.Group:
		switch any(zero).(type) {
		case api.GroupsResponse:
			return any(convertGroupToResponse(*v)).(T), nil
		}
	}

	return zero, errors.New("unsupported group conversion")
}

// 内部の詰め替え用補助関数
func convertGroupToResponse(group model.Group) api.GroupsResponse {
	resp := api.GroupsResponse{
		Id:      group.ID,
		Name:    group.Name,
		Members: make([]api.GroupMember, len(group.Members)),
	}

	for i, member := range group.Members {
		resp.Members[i] = api.GroupMember{
			Uid:  member.UID,
			Name: member.Name,
		}
	}

	return resp
}