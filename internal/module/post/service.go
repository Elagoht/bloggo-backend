package post

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"bloggo/internal/utils/readtime"
	"bloggo/internal/utils/schemas/responses"
	"fmt"
)

type PostService struct {
	repository     PostRepository
	bucket         bucket.Bucket
	imageValidator validatefile.FileValidator
	coverResizer   transformfile.FileTransformer
	permissions    permissions.Store
}

func NewPostService(
	repository PostRepository,
	bucket bucket.Bucket,
	imageValidator validatefile.FileValidator,
	coverResizer transformfile.FileTransformer,
	permissions permissions.Store,
) PostService {
	return PostService{
		repository,
		bucket,
		imageValidator,
		coverResizer,
		permissions,
	}
}

func (service *PostService) GetPostList() (
	[]models.ResponsePostCard,
	error,
) {
	return service.repository.GetPostList()
}

func (service *PostService) GetPostListPaginated(
	filters *models.RequestPostFilters,
) (*responses.PaginatedResponse[models.ResponsePostCard], error) {
	response, err := service.repository.GetPostListPaginated(filters)
	if err != nil {
		return nil, err
	}

	// Add avatar URL prefix to each post author if avatar exists
	for i := range response.Data {
		if response.Data[i].Author.Avatar != nil && *response.Data[i].Author.Avatar != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *response.Data[i].Author.Avatar)
			response.Data[i].Author.Avatar = &avatarPath
		}
	}

	return response, nil
}

func (service *PostService) GetPostById(
	id int64,
) (*models.ResponsePostDetails, error) {
	post, err := service.repository.GetPostById(id)
	if err != nil {
		return nil, err
	}

	// Add avatar URL prefix if avatar exists
	if post.Author.Avatar != nil && *post.Author.Avatar != "" {
		avatarPath := fmt.Sprintf("/uploads/avatar/%s", *post.Author.Avatar)
		post.Author.Avatar = &avatarPath
	}

	return post, nil
}

func (service *PostService) GetPostBySlug(
	slug string,
) (*models.ResponsePostDetails, error) {
	post, err := service.repository.GetPostGetByCurrentVersionSlug(slug)
	if err != nil {
		return nil, err
	}

	// Add avatar URL prefix if avatar exists
	if post.Author.Avatar != nil && *post.Author.Avatar != "" {
		avatarPath := fmt.Sprintf("/uploads/avatar/%s", *post.Author.Avatar)
		post.Author.Avatar = &avatarPath
	}

	return post, nil
}

func (service *PostService) CreatePostWithFirstVersion(
	model *models.RequestPostUpsert,
	userId int64,
) (*responses.ResponseCreated, error) {
	var filePath string

	// Handle cover file if provided
	if model.Cover != nil {
		coverFile, err := model.Cover.Open()
		if err != nil {
			return nil, err
		}
		defer coverFile.Close()

		filePath = cryptography.GenerateUniqueId() + ".webp"

		if err = service.imageValidator.Validate(coverFile, model.Cover); err != nil {
			return nil, err
		}

		transformedFile, err := service.coverResizer.Transform(coverFile)
		if err != nil {
			return nil, err
		}

		service.bucket.Save(transformedFile, filePath)
	}

	// Calculate read time
	content := ""
	if model.Content != nil {
		content = *model.Content
	}
	estimatedReadTime := readtime.EstimateReadTime(content)

	createdId, err := service.repository.CreatePost(model, filePath, estimatedReadTime, userId)
	if err != nil {
		// If cannot created and file was uploaded, delete it
		if filePath != "" {
			service.bucket.Delete(filePath)
		}
		return nil, err
	}

	return &responses.ResponseCreated{
		Id: createdId,
	}, nil
}

func (service *PostService) ListPostVersionsGetByPostId(
	id int64,
) (*models.ResponseVersionsOfPost, error) {
	response, err := service.repository.ListPostVersionsGetByPostId(id)
	if err != nil {
		return nil, err
	}

	// Add avatar URL prefix to original author if avatar exists
	if response.OriginalAuthor.Avatar != nil && *response.OriginalAuthor.Avatar != "" {
		avatarPath := fmt.Sprintf("/uploads/avatar/%s", *response.OriginalAuthor.Avatar)
		response.OriginalAuthor.Avatar = &avatarPath
	}

	// Add avatar URL prefix to each version author if avatar exists
	for i := range response.Versions {
		if response.Versions[i].VersionAuthor.Avatar != nil && *response.Versions[i].VersionAuthor.Avatar != "" {
			avatarPath := fmt.Sprintf("/uploads/avatar/%s", *response.Versions[i].VersionAuthor.Avatar)
			response.Versions[i].VersionAuthor.Avatar = &avatarPath
		}
	}

	return response, nil
}

func (service *PostService) GetPostVersionById(
	postId int64,
	versionId int64,
) (*models.ResponseVersionDetailsOfPost, error) {
	version, err := service.repository.GetPostVersionById(postId, versionId)
	if err != nil {
		return nil, err
	}

	// Add avatar URL prefix if avatar exists
	if version.VersionAuthor.Avatar != nil && *version.VersionAuthor.Avatar != "" {
		avatarPath := fmt.Sprintf("/uploads/avatar/%s", *version.VersionAuthor.Avatar)
		version.VersionAuthor.Avatar = &avatarPath
	}

	return version, nil
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

	// Calculate read time if content is being updated
	var readTime *int
	if model.Content != nil {
		calculatedReadTime := readtime.EstimateReadTime(*model.Content)
		readTime = &calculatedReadTime
	}

	if err := service.repository.UpdateVersionById(
		postId,
		versionId,
		userId,
		model,
		filePath,
		readTime,
	); err != nil {
		// If cannot created, delete newly uploaded file
		if filePath != nil {
			service.bucket.Delete(*filePath)
		}
		return err
	}

	return nil
}

func (service *PostService) SubmitVersionForReview(
	postId int64,
	versionId int64,
	userId int64,
) error {
	// Check if the owner of version is same as requester
	versionCreator, versionStatus, err :=
		service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return err
	}

	// Users can only submit their own versions
	if versionCreator != userId {
		return apierrors.ErrForbidden
	}

	// Only draft versions can be submitted
	if versionStatus != models.STATUS_DRAFT {
		return apierrors.ErrPreconditionFailed
	}

	// Update version status to pending (submitted for review)
	return service.repository.UpdateVersionStatus(
		versionId,
		models.STATUS_PENDING,
		userId,
	)
}

func (service *PostService) ApproveVersion(
	postId int64,
	versionId int64,
	userId int64,
	note *string,
) error {
	// Check if version exists and get current status
	_, versionStatus, err := service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return err
	}

	// Only pending versions can be approved, drafts and published cannot be approved
	if versionStatus == models.STATUS_DRAFT || versionStatus == models.STATUS_PUBLISHED {
		return apierrors.ErrPreconditionFailed
	}

	// Update version status to approved
	return service.repository.UpdateVersionStatusWithNote(
		versionId,
		models.STATUS_APPROVED,
		userId,
		note,
	)
}

func (service *PostService) RejectVersion(
	postId int64,
	versionId int64,
	userId int64,
	note *string,
) error {
	// Check if version exists and get current status
	_, versionStatus, err := service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return err
	}

	// Only pending or approved versions can be rejected, drafts and published cannot be rejected
	if versionStatus == models.STATUS_DRAFT || versionStatus == models.STATUS_PUBLISHED {
		return apierrors.ErrPreconditionFailed
	}

	// Update version status to rejected
	return service.repository.UpdateVersionStatusWithNote(
		versionId,
		models.STATUS_REJECTED,
		userId,
		note,
	)
}

func (service *PostService) DeleteVersionById(
	postId int64,
	versionId int64,
	userId int64,
	roleId int64,
) error {
	// Get version details including creator and status
	versionCreator, versionStatus, err := service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return err
	}

	// Check if user has editor permissions (can delete any version)
	hasEditorPermission := service.permissions.HasPermission(roleId, "post:delete")

	// Check if user owns the version (can delete own versions with restrictions)
	isOwner := versionCreator == userId

	if !hasEditorPermission && !isOwner {
		return apierrors.ErrForbidden
	}

	// If user is owner but not editor, check status restrictions
	if isOwner && !hasEditorPermission {
		// Authors can only delete draft, pending, or rejected versions
		if versionStatus != models.STATUS_DRAFT &&
			versionStatus != models.STATUS_PENDING &&
			versionStatus != models.STATUS_REJECTED {
			return apierrors.ErrPreconditionFailed
		}
	}

	// Get the cover image path before deletion
	coverImagePath, err := service.repository.GetVersionCoverImage(versionId)
	if err != nil {
		return err
	}

	// Check if this version is currently published
	isCurrentlyPublished, err := service.repository.IsVersionCurrentlyPublished(versionId)
	if err != nil {
		return err
	}

	// If it's currently published, set the post's current_version_id to NULL
	if isCurrentlyPublished {
		if err := service.repository.SetPostCurrentVersionToNull(versionId); err != nil {
			return err
		}
	}

	// Perform soft delete
	if err := service.repository.SoftDeleteVersionById(versionId); err != nil {
		return err
	}

	// Check if the cover image is still referenced by other versions
	if coverImagePath != nil {
		isImageStillReferenced, err :=
			service.repository.IsImageReferencedByOtherVersions(
				*coverImagePath,
				versionId,
			)
		if err != nil {
			// Log error but don't fail the deletion
			return nil
		}

		// If image is not referenced by any other version, delete it from storage
		if !isImageStillReferenced {
			service.bucket.Delete(*coverImagePath)
		}
	}

	return nil
}

func (service *PostService) TrackView(model *models.RequestTrackView) error {
	return service.repository.TrackView(model.PostId, model.UserAgent)
}
