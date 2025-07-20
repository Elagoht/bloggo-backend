package user

import (
	"bloggo/internal/module/user/models"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
)

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) UserService {
	return UserService{
		repository,
	}
}

func (service *UserService) GetUsers(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseUserCard, error) {
	return service.repository.GetUsers(paginate, search)
}

func (service *UserService) GetUserById(
	id int,
) (*models.ResponseUserDetails, error) {
	return service.repository.GetUserById(id)
}
