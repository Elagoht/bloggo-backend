package post

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"mime/multipart"
)

type PostService struct {
	repository     PostRepository
	bucket         bucket.Bucket
	imageValidator validatefile.FileValidator
	coverResizer   transformfile.FileTransformer
}

func NewPostService(
	repository PostRepository,
	bucket bucket.Bucket,
	imageValidator validatefile.FileValidator,
	coverResizer transformfile.FileTransformer,
) PostService {
	return PostService{
		repository,
		bucket,
		imageValidator,
		coverResizer,
	}
}

func (service *PostService) GetPostList() (
	[]models.ResponsePostCard,
	error,
) {
	return service.repository.GetPostList()
}

func (service *PostService) GetPostById(
	id int64,
) (*models.ResponsePostDetails, error) {
	return service.repository.GetPostById(id)
}

func (service *PostService) CreatePostWithFirstVersion(
	model *models.RequestPostUpsert,
	cover *multipart.FileHeader,
	userId int64,
) (*models.ResponsePostCreated, error) {
	coverFile, err := cover.Open()
	if err != nil {
		return nil, err
	}
	defer coverFile.Close()

	filepath := cryptography.GenerateUniqueId() + ".webp"

	if err = service.imageValidator.Validate(coverFile, cover); err != nil {
		return nil, err
	}

	transformedFile, err := service.coverResizer.Transform(coverFile)
	if err != nil {
		return nil, err
	}

	service.bucket.Save(transformedFile, filepath)

	createdId, err := service.repository.CreatePost(model, filepath, userId)
	if err != nil {
		// If cannot created, delete newly uploaded file
		service.bucket.Delete(filepath)
		return nil, err
	}

	return &models.ResponsePostCreated{
		Id: createdId,
	}, nil
}
