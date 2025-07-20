package user

import (
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
	id, ok := handlers.GetParam[int](writer, request, "id")
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
