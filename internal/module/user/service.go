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

func (service *UserService) UserCreate(
	model *models.RequestUserCreate,
) (*models.ResponseUserCreated, error) {
	processed, err := model.HashUserPassphrase()
	if err != nil {
		return nil, err
	}

	id, err := service.repository.UserCreate(processed)
	if err != nil {
		return nil, err
	}

	return &models.ResponseUserCreated{
		Id: id,
	}, nil
}
