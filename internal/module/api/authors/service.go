package authors

import (
	"bloggo/internal/module/api/authors/models"
)

type AuthorsAPIService struct {
	repository AuthorsAPIRepository
}

func NewAuthorsAPIService(repository AuthorsAPIRepository) AuthorsAPIService {
	return AuthorsAPIService{repository}
}

func (service *AuthorsAPIService) GetAllAuthors() (*models.APIAuthorsResponse, error) {
	return service.repository.GetAllAuthors()
}

func (service *AuthorsAPIService) GetAuthorById(id int64) (*models.APIAuthorDetails, error) {
	return service.repository.GetAuthorById(id)
}
