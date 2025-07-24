package user

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/module/user/models"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"fmt"
	"log"
	"mime/multipart"
)

type UserService struct {
	repository     UserRepository
	bucket         bucket.Bucket
	imageValidator validatefile.FileValidator
	avatarResizer  transformfile.FileTransformer
}

func NewUserService(
	repository UserRepository,
	bucket bucket.Bucket,
	imageValidator validatefile.FileValidator,
	avatarResizer transformfile.FileTransformer,
) UserService {
	return UserService{
		repository,
		bucket,
		imageValidator,
		avatarResizer,
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
	// Check if file is an image
	if err := service.imageValidator.Validate(file, header); err != nil {
		return err
	}

	// Resize and format image as webp
	converted, err := service.avatarResizer.Transform(file)
	if err != nil {
		return err
	}

	// Create a unique name but related to user via a seperator.
	fileName := service.createUserRelatedUUID(userID) + ".webp"

	// Save new avatar
	if err := service.bucket.Save(converted, fileName); err != nil {
		return fmt.Errorf("failed to save avatar: %w", err)
	}

	// Delete old avatar
	if err := service.bucket.DeleteMatching(
		fmt.Sprintf("%d_*.webp", userID),
		fileName, // Protect new image by blacklisting it
	); err != nil {
		return fmt.Errorf("failed to delete old avatars: %w", err)
	}

	// Update database
	if err := service.repository.UpdateAvatarById(userID, fileName); err != nil {
		if deleteErr := service.bucket.Delete(fileName); deleteErr != nil {
			log.Printf("Failed to delete avatar after db error: %v", deleteErr)
		}
		return fmt.Errorf("failed to update avatar in database: %w", err)
	}

	return nil
}

func (service *UserService) createUserRelatedUUID(
	userID int64,
) string {
	uuid := cryptography.GenerateUniqueId()
	return fmt.Sprintf("%d_%s", userID, uuid)
}
