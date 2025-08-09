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
	userId, ok := handlers.GetContextValue[int64](writer, request, "userId")
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
	body, ok := handlers.BindAndValidate[*models.RequestUserCreate](writer, request)
	if !ok {
		return
	}

	created, err := handler.service.UserCreate(body)
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
	userId, ok := handlers.GetContextValue[int64](writer, request, "userId")
	if !ok {
		return
	}

	file, fileHeader, ok := handlers.GetFormFile(writer, request, "avatar", 10<<20)
	if !ok {
		return
	}
	defer file.Close()

	err := handler.service.UpdateAvatarById(userId, file, fileHeader)
	if err != nil {
		apierrors.MapErrors(err, writer, nil)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}
