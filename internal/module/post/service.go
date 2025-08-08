package post

import (
	"bloggo/internal/module/post/models"
)

type PostService struct {
	repository PostRepository
}

func NewPostService(repository PostRepository) PostService {
	return PostService{
		repository,
	}
}

func (service *PostService) GetPostList() (
	[]models.ResponsePostCard,
	error,
) {
	return service.repository.GetPostList()
}

func (service *PostService) GetPostBySlug(
	slug string,
) (*models.ResponsePostDetails, error) {
	return service.repository.GetPostBySlug(slug)
}
