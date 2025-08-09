package post

import (
	"bloggo/internal/module/post/models"
	"mime/multipart"
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

func (service *PostService) CreatePostWithFirstVersion(
	model *models.RequestPostUpsert,
	cover *multipart.FileHeader,
	userId int64,
) (*models.ResponsePostCreated, error) {
	// TODO: Upload file

	createdId, err := service.repository.CreatePost(model, "", userId)
	if err != nil {
		return nil, err
	}

	return &models.ResponsePostCreated{
		Id: createdId,
	}, nil
}
