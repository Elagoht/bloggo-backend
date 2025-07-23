package user

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/module/user/models"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"mime/multipart"
	"strconv"
)

type UserService struct {
	repository UserRepository
	bucket     bucket.Bucket
}

func NewUserService(
	repository UserRepository,
	bucket bucket.Bucket,
) UserService {
	return UserService{
		repository,
		bucket,
	}
}

func (service *UserService) GetUsers(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseUserCard, error) {
	return service.repository.GetUsers(paginate, search)
}

func (service *UserService) GetUserById(
	id int64,
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

func (service *UserService) UpdateAvatarById(
	userID int64,
	file multipart.File,
	header *multipart.FileHeader,
) error {
	// Create a unique name but related to user via a seperator.
	fileName, err := service.createUserRelatedUUID(userID)
	if err != nil {
		return err
	}

	fileNameWithExtension := fileName + ".png"

	// Upsert avatar image by userID to storage
	if err := service.bucket.SaveMultiPart(
		file,
		fileNameWithExtension,
	); err != nil {
		return err
	}

	// Update new image URL on database
	if err := service.repository.UpdateAvatarById(
		userID,
		fileNameWithExtension,
	); err != nil {
		return err
	}

	return nil
}

func (service *UserService) createUserRelatedUUID(
	userID int64,
) (string, error) {
	uuid, err := cryptography.GenerateUniqueId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(userID, 10) + "_" + uuid, nil
}
