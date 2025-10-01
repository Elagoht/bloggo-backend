package post

import (
	"bloggo/internal/infrastructure/bucket"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/ai"
	aimodels "bloggo/internal/module/ai/models"
	"bloggo/internal/module/audit"
	auditmodels "bloggo/internal/module/audit/models"
	"bloggo/internal/module/post/models"
	"bloggo/internal/utils/apierrors"
	"bloggo/internal/utils/cryptography"
	"bloggo/internal/utils/file/transformfile"
	"bloggo/internal/utils/file/validatefile"
	"bloggo/internal/utils/readtime"
	"bloggo/internal/utils/schemas/responses"
	"bloggo/internal/utils/validate"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

type PostService struct {
	repository     PostRepository
	bucket         bucket.Bucket
	imageValidator validatefile.FileValidator
	coverResizer   transformfile.FileTransformer
	permissions    permissions.Store
	aiService      ai.AIService
	cache          *GenerativeFillCache
}

type CacheEntry struct {
	data      *aimodels.ResponseGenerativeFill
	timestamp time.Time
}

type GenerativeFillCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheEntry
}

func NewGenerativeFillCache() *GenerativeFillCache {
	return &GenerativeFillCache{
		cache: make(map[string]*CacheEntry),
	}
}

func (c *GenerativeFillCache) Get(key string) (*aimodels.ResponseGenerativeFill, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}

	// Check if cache entry is expired (1 minute)
	if time.Since(entry.timestamp) > time.Minute {
		delete(c.cache, key)
		return nil, false
	}

	return entry.data, true
}

func (c *GenerativeFillCache) Set(key string, data *aimodels.ResponseGenerativeFill) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &CacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
}

func NewPostService(
	repository PostRepository,
	bucket bucket.Bucket,
	imageValidator validatefile.FileValidator,
	coverResizer transformfile.FileTransformer,
	permissions permissions.Store,
) PostService {
	return PostService{
		repository:     repository,
		bucket:         bucket,
		imageValidator: imageValidator,
		coverResizer:   coverResizer,
		permissions:    permissions,
		aiService:      ai.NewAIService(),
		cache:          NewGenerativeFillCache(),
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
			avatarPath := fmt.Sprintf(
				"/uploads/avatar/%s",
				*response.Data[i].Author.Avatar,
			)
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

	// Get tags for this post
	tags, err := service.repository.GetPostTags(id)
	if err != nil {
		// Don't fail if tags can't be retrieved, just return empty array
		tags = []models.TagCard{}
	}
	post.Tags = tags

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

	// Get tags for this post
	tags, err := service.repository.GetPostTags(post.PostId)
	if err != nil {
		// Don't fail if tags can't be retrieved, just return empty array
		tags = []models.TagCard{}
	}
	post.Tags = tags

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

	// Log post creation audit
	audit.LogPostAction(&userId, createdId, auditmodels.ActionPostCreated)

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

func (service *PostService) CreateVersionFromSpecificVersion(
	versionId int64,
	userId int64,
) (*responses.ResponseCreated, error) {
	createdId, err := service.repository.CreateVersionFromSpecificVersion(versionId, userId)
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
		service.repository.GetVersionCreatorAndStatus(versionId)
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

	// Validate the version content before allowing submission
	if err := service.validateVersionForSubmission(postId, versionId); err != nil {
		return err
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
) (*models.ResponseVersionDeleted, error) {
	// Get version details including creator and status
	versionCreator, versionStatus, err := service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return nil, err
	}

	// Check if user has editor permissions (can delete any version)
	hasEditorPermission := service.permissions.HasPermission(roleId, "post:delete")

	// Check if user owns the version (can delete own versions with restrictions)
	isOwner := versionCreator == userId

	if !hasEditorPermission && !isOwner {
		return nil, apierrors.ErrForbidden
	}

	// If user is owner but not editor, check status restrictions
	if isOwner && !hasEditorPermission {
		// Authors can only delete draft, pending, or rejected versions
		if versionStatus != models.STATUS_DRAFT &&
			versionStatus != models.STATUS_PENDING &&
			versionStatus != models.STATUS_REJECTED {
			return nil, apierrors.ErrPreconditionFailed
		}
	}

	// Check if this is the last version of the post
	versionCount, err := service.repository.CountPostVersions(postId)
	if err != nil {
		return nil, err
	}

	// If this is the last version, delete the entire post instead
	if versionCount == 1 {
		// Store cover photo paths before deleting post
		coverPaths, err := service.repository.GetAllRelatedCovers(postId)
		if err != nil {
			return nil, err
		}

		// Delete the entire post (soft delete)
		if err := service.repository.SoftDeletePostById(postId); err != nil {
			return nil, err
		}

		// Delete all related images
		for _, path := range coverPaths {
			service.bucket.Delete(path)
		}

		// Log post deletion audit
		audit.LogPostAction(&userId, postId, auditmodels.ActionPostDeleted)

		return &models.ResponseVersionDeleted{PostDeleted: true}, nil
	}

	// Get the cover image path before deletion
	coverImagePath, err := service.repository.GetVersionCoverImage(versionId)
	if err != nil {
		return nil, err
	}

	// Check if this version is currently published
	isCurrentlyPublished, err := service.repository.IsVersionCurrentlyPublished(versionId)
	if err != nil {
		return nil, err
	}

	// If it's currently published, set the post's current_version_id to NULL
	if isCurrentlyPublished {
		if err := service.repository.SetPostCurrentVersionToNull(versionId); err != nil {
			return nil, err
		}
	}

	// Perform soft delete
	if err := service.repository.SoftDeleteVersionById(versionId); err != nil {
		return nil, err
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
			return &models.ResponseVersionDeleted{PostDeleted: false}, nil
		}

		// If image is not referenced by any other version, delete it from storage
		if !isImageStillReferenced {
			service.bucket.Delete(*coverImagePath)
		}
	}

	// Log version deletion audit
	audit.LogVersionAction(&userId, versionId, auditmodels.ActionVersionDeleted, nil)

	return &models.ResponseVersionDeleted{PostDeleted: false}, nil
}

func (service *PostService) PublishVersion(
	postId int64,
	versionId int64,
	userId int64,
	roleId int64,
) error {
	// Check if user has publish permission
	hasPublishPermission := service.permissions.HasPermission(roleId, "post:publish")
	if !hasPublishPermission {
		return apierrors.ErrForbidden
	}

	// Check if version exists and get current status
	_, versionStatus, err := service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return err
	}

	// Only approved versions can be published
	if versionStatus != models.STATUS_APPROVED {
		return apierrors.ErrPreconditionFailed
	}

	// Check if the version's category is deleted
	categoryIsDeleted, err := service.repository.CheckIfVersionCategoryIsDeleted(versionId)
	if err != nil {
		return err
	}

	if categoryIsDeleted {
		// Return 428 Precondition Required status to indicate category needs to be updated
		return apierrors.NewAPIError(
			"This version's category has been deleted. Please select a new category before publishing.",
			apierrors.ErrPreconditionRequired,
		)
	}

	// Get the slug of the version being published
	slug, err := service.repository.GetVersionSlug(versionId)
	if err != nil {
		return err
	}

	// Check if there's already a published version with the same slug
	existingPublished, err := service.repository.GetPublishedVersionBySlug(slug)
	if err == nil && existingPublished != nil {
		// Unpublish the existing version (set it back to approved status)
		if err := service.repository.UnpublishVersionBySlug(slug); err != nil {
			return err
		}

		// Clear the current_version_id from the post that was using the old published version
		if err := service.repository.SetPostCurrentVersionToNull(existingPublished.Id); err != nil {
			return err
		}
	}

	// Update version status to published
	if err := service.repository.UpdateVersionStatus(
		versionId,
		models.STATUS_PUBLISHED,
		userId,
	); err != nil {
		return err
	}

	// Log version publication audit
	audit.LogVersionAction(&userId, versionId, auditmodels.ActionVersionPublished, nil)

	// Set this version as the current published version for the post
	return service.repository.SetCurrentVersionForPost(postId, versionId)
}

func (service *PostService) TrackView(model *models.RequestTrackView) error {
	return service.repository.TrackView(model.PostId, model.UserAgent)
}

// validateVersionForSubmission validates that a post version has all required fields
// populated before it can be submitted for review
func (service *PostService) validateVersionForSubmission(postId int64, versionId int64) error {
	// Get the version details to validate
	version, err := service.repository.GetPostVersionById(postId, versionId)
	if err != nil {
		return err
	}

	// Create validation struct from version data
	validationData := models.PostSubmissionValidation{}

	// Map version fields to validation struct, handling nil pointers
	if version.Title != nil {
		validationData.Title = strings.TrimSpace(*version.Title)
	}
	if version.Content != nil {
		validationData.Content = strings.TrimSpace(*version.Content)
	}
	if version.Description != nil {
		validationData.Description = strings.TrimSpace(*version.Description)
	}
	if version.Spot != nil {
		validationData.Spot = strings.TrimSpace(*version.Spot)
	}
	if version.Category.Id != nil {
		if categoryId, err := strconv.ParseInt(*version.Category.Id, 10, 64); err == nil {
			validationData.CategoryId = categoryId
		}
	}

	// Use the validator to validate the struct
	validatorInstance := validate.GetValidator()
	if err := validatorInstance.Struct(validationData); err != nil {
		// Convert validator errors to API errors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errors []apierrors.ValidationError
			for _, validationError := range validationErrors {
				errors = append(errors, apierrors.ValidationError{
					Field:   validationError.Field(),
					Message: service.getValidationErrorMessage(validationError),
				})
			}
			return apierrors.NewValidationAPIError(errors)
		}
		return err
	}

	return nil
}

// getValidationErrorMessage converts validation error tags to human-readable messages
func (service *PostService) getValidationErrorMessage(
	fieldError validator.FieldError,
) string {
	switch fieldError.Tag() {
	case "required":
		return fieldError.Field() + " is required"
	case "min":
		return fieldError.Field() + " must be at least " + fieldError.Param() + " characters"
	case "max":
		return fieldError.Field() + " must be " + fieldError.Param() + " characters or less"
	default:
		return fieldError.Field() + " is invalid"
	}
}

func (service *PostService) GenerativeFill(
	postId,
	versionId int64,
	availableCategories []string,
) (*aimodels.ResponseGenerativeFill, error) {
	// Get the version details
	version, err := service.repository.GetPostVersionById(postId, versionId)
	if err != nil {
		return nil, err
	}

	// Check if content exists and has minimum 1000 characters
	if version.Content == nil || len(*version.Content) < 1000 {
		return nil, apierrors.NewAPIError(
			"Content must be at least 1000 characters long",
			apierrors.ErrBadRequest,
		)
	}

	// Create cache key based on content hash
	cacheKey := fmt.Sprintf(
		"%d-%d-%s",
		postId,
		versionId,
		cryptography.HashString(*version.Content),
	)

	// Check cache first
	if cached, found := service.cache.Get(cacheKey); found {
		return cached, nil
	}

	// Generate AI metadata
	result, err := service.aiService.GenerateContentMetadata(*version.Content, availableCategories)
	if err != nil {
		return nil, err
	}

	// Cache the result for 1 minute
	service.cache.Set(cacheKey, result)

	return result, nil
}

func (service *PostService) AssignTagsToPost(
	postId int64,
	tagIds []int64,
	userRoleId int64,
) error {
	// Check if user has permission to assign tags
	hasPermission := service.permissions.HasPermission(userRoleId, "tag:assign")
	if !hasPermission {
		return apierrors.ErrForbidden
	}

	return service.repository.AssignTagsToPost(postId, tagIds)
}

func (service *PostService) UpdateVersionCategory(
	postId int64,
	versionId int64,
	categoryId int64,
	userId int64,
	roleId int64,
) error {
	// Check if user has publish permission (required to update approved version category)
	hasPublishPermission := service.permissions.HasPermission(roleId, "post:publish")
	if !hasPublishPermission {
		return apierrors.ErrForbidden
	}

	// Check if version exists and get current status
	_, versionStatus, err := service.repository.GetVersionCreatorAndStatus(versionId)
	if err != nil {
		return err
	}

	// Only approved versions should need category updates (this is the special case)
	if versionStatus != models.STATUS_APPROVED {
		return apierrors.ErrPreconditionFailed
	}

	// Check if the version's category is actually deleted (prevent abuse)
	categoryIsDeleted, err := service.repository.CheckIfVersionCategoryIsDeleted(versionId)
	if err != nil {
		return err
	}

	if !categoryIsDeleted {
		return apierrors.NewAPIError(
			"The version's category is not deleted. This endpoint is only for updating categories when the original has been deleted.",
			apierrors.ErrPreconditionFailed,
		)
	}

	// Update the category
	return service.repository.UpdateVersionCategoryOnly(versionId, categoryId)
}
