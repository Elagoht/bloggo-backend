package user

import (
	"bloggo/internal/module/user/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) UserHandler {
	return UserHandler{
		service,
	}
}

func (handler *UserHandler) GetSelf(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	user, err := handler.service.GetUserById(userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(user)
}

func (handler *UserHandler) GetUsers(
	writer http.ResponseWriter,
	request *http.Request,
) {
	paginate, ok := pagination.GetPaginationOptions(writer, request, []string{
		"name", "created_at", "updated_at", "last_login",
	})
	if !ok {
		return
	}

	// Map API field names to qualified SQL column names
	if paginate.OrderBy != nil && *paginate.OrderBy == "name" {
		qualifiedName := "users.name"
		paginate.OrderBy = &qualifiedName
	}

	search, ok := filter.GetSearchOptions(writer, request)
	if !ok {
		return
	}

	users, err := handler.service.GetUsers(paginate, search)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(users)
}

func (handler *UserHandler) GetUserById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	user, err := handler.service.GetUserById(id)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(user)
}

func (handler *UserHandler) UserCreate(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestUserCreate](writer, request)
	if !ok {
		return
	}

	created, err := handler.service.UserCreate(body, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	json.NewEncoder(writer).Encode(created)
}

func (handler *UserHandler) UpdateSelfAvatar(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	file, fileHeader, ok := handlers.GetFormFile(writer, request, "avatar", 10<<20)
	if !ok {
		return
	}
	defer file.Close()

	avatarPath, err := handler.service.UpdateAvatarById(userId, file, fileHeader, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	response := models.ResponseAvatarUpdate{Avatar: avatarPath}
	json.NewEncoder(writer).Encode(response)
}

func (handler *UserHandler) DeleteUserAvatar(
	writer http.ResponseWriter,
	request *http.Request,
) {
	deleterId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	err := handler.service.DeleteAvatarById(id, deleterId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *UserHandler) DeleteSelfAvatar(
	writer http.ResponseWriter,
	request *http.Request,
) {
	userId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	err := handler.service.DeleteAvatarById(userId, userId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *UserHandler) UpdateUserById(
	writer http.ResponseWriter,
	request *http.Request,
) {
	updaterId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestUserUpdate](writer, request)
	if !ok {
		return
	}

	err := handler.service.UpdateUserById(id, body, updaterId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *UserHandler) UpdateUserAvatar(
	writer http.ResponseWriter,
	request *http.Request,
) {
	updaterId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	file, fileHeader, ok := handlers.GetFormFile(writer, request, "avatar", 10<<20)
	if !ok {
		return
	}
	defer file.Close()

	avatarPath, err := handler.service.UpdateAvatarById(id, file, fileHeader, updaterId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	response := models.ResponseAvatarUpdate{Avatar: avatarPath}
	json.NewEncoder(writer).Encode(response)
}

func (handler *UserHandler) AssignRole(
	writer http.ResponseWriter,
	request *http.Request,
) {
	assignerId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestUserAssignRole](writer, request)
	if !ok {
		return
	}

	err := handler.service.AssignRole(id, body, assignerId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *UserHandler) ChangePassword(
	writer http.ResponseWriter,
	request *http.Request,
) {
	changerId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	body, ok := handlers.BindAndValidate[*models.RequestUserChangePassword](writer, request)
	if !ok {
		return
	}

	err := handler.service.ChangePassword(id, body, changerId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func (handler *UserHandler) DeleteUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	deleterId, ok := handlers.GetContextValue[int64](writer, request, handlers.TokenUserId)
	if !ok {
		return
	}

	id, ok := handlers.GetParam[int64](writer, request, "id")
	if !ok {
		return
	}

	err := handler.service.DeleteUser(id, deleterId)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
