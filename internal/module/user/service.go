package user

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/module/user/models"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/pagination"
	"bloggo/internal/utils/schemas/responses"
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
) (*responses.PaginatedResponse[models.ResponseUserCard], error) {
	// Get the users data
	users, err := service.repository.GetUsers(paginate, search)
	if err != nil {
		return nil, err
	}

	// Add avatar prefix to each user
	for index := range users {
		if users[index].Avatar != nil && *users[index].Avatar != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *users[index].Avatar)
			users[index].Avatar = &avatarPath
		}
	}

	// Get the total count with same filters
	total, err := service.repository.GetUsersCount(search)
	if err != nil {
		return nil, err
	}

	// Get page number (default to 1)
	page := 1
	if paginate.Page != nil {
		page = *paginate.Page
	}

	// Get take value (default to 10)
	take := 10
	if paginate.Take != nil {
		take = *paginate.Take
	}

	return &responses.PaginatedResponse[models.ResponseUserCard]{
		Data:  users,
		Page:  page,
		Take:  take,
		Total: total,
	}, nil
}

func (service *UserService) GetUserById(
	id int64,
) (*models.ResponseUserDetails, error) {
	user, err := service.repository.GetUserById(id)
	if err != nil {
		return nil, err
	}

	if user.Avatar != nil && *user.Avatar != "" {
		*user.Avatar = fmt.Sprintf("/uploads/avatar/%s", *user.Avatar)
	}
	return user, nil
}

func (service *UserService) UserCreate(
	model *models.RequestUserCreate,
) (*responses.ResponseCreated, error) {
	processed, err := model.HashUserPassphrase()
	if err != nil {
		return nil, err
	}

	id, err := service.repository.UserCreate(processed)
	if err != nil {
		return nil, err
	}

	return &responses.ResponseCreated{
		Id: id,
	}, nil
}

func (service *UserService) UpdateAvatarById(
	userId int64,
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
	imageId := service.createUserRelatedUUID(userId)
	fileName := imageId + ".webp"

	// Save new avatar
	if err := service.bucket.Save(converted, fileName); err != nil {
		return fmt.Errorf("failed to save avatar: %w", err)
	}

	// Delete old avatar
	if err := service.bucket.DeleteMatching(
		fmt.Sprintf("%d_*.webp", userId),
		fileName, // Protect new image by blacklisting it
	); err != nil {
		return fmt.Errorf("failed to delete old avatars: %w", err)
	}

	// Update database
	if err := service.repository.UpdateAvatarById(userId, imageId); err != nil {
		if deleteErr := service.bucket.Delete(fileName); deleteErr != nil {
			log.Printf("Failed to delete avatar after db error: %v", deleteErr)
		}
		return fmt.Errorf("failed to update avatar in database: %w", err)
	}

	return nil
}

func (service *UserService) UpdateUserById(
	userId int64,
	model *models.RequestUserUpdate,
) error {
	return service.repository.UpdateUserById(userId, model)
}

func (service *UserService) AssignRole(
	userId int64,
	model *models.RequestUserAssignRole,
) error {
	return service.repository.AssignRole(userId, model.RoleId)
}

func (service *UserService) DeleteUser(userId int64) error {
	return service.repository.DeleteUser(userId)
}

func (service *UserService) UpdateLastLogin(userId int64) error {
	return service.repository.UpdateLastLogin(userId)
}

func (service *UserService) ChangePassword(
	userId int64,
	model *models.RequestUserChangePassword,
) error {
	hashedPassword, err := model.HashNewPassword()
	if err != nil {
		return err
	}

	return service.repository.UpdatePasswordById(userId, hashedPassword)
}

func (service *UserService) DeleteAvatarById(userId int64) error {
	// Delete all avatar files for this user
	if err := service.bucket.DeleteMatching(
		fmt.Sprintf("%d_*.webp", userId),
		"", // No files to protect
	); err != nil {
		return fmt.Errorf("failed to delete avatar files: %w", err)
	}

	// Update database to remove avatar reference
	if err := service.repository.UpdateAvatarById(userId, ""); err != nil {
		return fmt.Errorf("failed to update avatar in database: %w", err)
	}

	return nil
}

func (service *UserService) createUserRelatedUUID(
	userId int64,
) string {
	uuid := cryptography.GenerateUniqueId()
	return fmt.Sprintf("%d_%s", userId, uuid)
}
