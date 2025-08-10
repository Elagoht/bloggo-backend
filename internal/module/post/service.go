package post

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"bloggo/internal/utils/schemas/responses"
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
	userId int64,
) (*responses.ResponseCreated, error) {
	coverFile, err := model.Cover.Open()
	if err != nil {
		return nil, err
	}
	defer coverFile.Close()

	filePath := cryptography.GenerateUniqueId() + ".webp"

	if err = service.imageValidator.Validate(coverFile, model.Cover); err != nil {
		return nil, err
	}

	transformedFile, err := service.coverResizer.Transform(coverFile)
	if err != nil {
		return nil, err
	}

	service.bucket.Save(transformedFile, filePath)

	createdId, err := service.repository.CreatePost(model, filePath, userId)
	if err != nil {
		// If cannot created, delete newly uploaded file
		service.bucket.Delete(filePath)
		return nil, err
	}

	return &responses.ResponseCreated{
		Id: createdId,
	}, nil
}

func (service *PostService) ListPostVersionsGetByPostId(
	id int64,
) (*models.ResponseVersionsOfPost, error) {
	return service.repository.ListPostVersionsGetByPostId(id)
}

func (service *PostService) GetPostVersionById(
	postId int64,
	versionId int64,
) (*models.ResponseVersionDetailsOfPost, error) {
	return service.repository.GetPostVersionById(postId, versionId)
}

func (service *PostService) DeletePostById(
	id int64,
) error {
	// Store cover photo paths before deleting post
	coverPaths, err := service.repository.GetAllRelatedCovers(id)
	if err != nil {
		return err
	}

	if err := service.repository.SoftDeletePostById(id); err != nil {
		return err
	}

	// If delete is succeed, delete photos
	for _, path := range coverPaths {
		service.bucket.Delete(path)
	}

	return nil
}

func (service *PostService) CreateVersionFromLatest(
	id int64,
	userId int64,
) (*responses.ResponseCreated, error) {
	createdId, err := service.repository.CreateVersionFromLatest(id, userId)
	if err != nil {
		return nil, err
	}

	return &responses.ResponseCreated{
		Id: createdId,
	}, nil
}

func (service *PostService) UpdateUnsubmittedOwnVersion(
	postId int64,
	versionId int64,
	userId int64,
	model *models.RequestPostUpsert,
) error {
	// 1. Check if the owner of version ismn same as requester
	versionCreator, versionStatus, err :=
		service.repository.GetVersionCreatorAndStatus(postId)
	if err != nil {
		return err
	}

	// Users can only edit their own versions
	if versionCreator != userId {
		return apierrors.ErrForbidden
	}

	// Only draft (unsubmitted versions can be edited)
	if versionStatus != models.STATUS_DRAFT {
		return apierrors.ErrPreconditionFailed
	}

	// If a new cover photo uploaded, validate, transform and save it
	var filePath *string
	if model.Cover != nil {
		coverFile, err := model.Cover.Open()
		if err != nil {
			return err
		}
		defer coverFile.Close()

		filepathName := cryptography.GenerateUniqueId() + ".webp"
		filePath = &filepathName

		if err = service.imageValidator.Validate(
			coverFile,
			model.Cover,
		); err != nil {
			return err
		}

		transformedFile, err := service.coverResizer.Transform(coverFile)
		if err != nil {
			return err
		}

		service.bucket.Save(transformedFile, *filePath)
	}

	if err := service.repository.UpdateVersionById(
		postId,
		versionId,
		userId,
		model,
		filePath,
	); err != nil {
		// If cannot created, delete newly uploaded file
		if filePath != nil {
			service.bucket.Delete(*filePath)
		}
		return err
	}

	return nil
}
