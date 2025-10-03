package webhook

import (
	"bloggo/internal/db"
	"bloggo/internal/infrastructure/permissions"
	"bloggo/internal/module/webhook/models"
	"sync"
	"time"
)

var (
	globalWebhookService WebhookService
	once                 sync.Once
)

// GetGlobalWebhookService returns the global webhook service instance
func GetGlobalWebhookService() WebhookService {
	once.Do(func() {
		database := db.Get()
		permissionStore := permissions.Get()
		repository := NewWebhookRepository(database)
		globalWebhookService = NewWebhookService(repository, permissionStore)
	})
	return globalWebhookService
}

// Helper functions for easy webhook triggering

// TriggerPostCreated fires a webhook for post creation
func TriggerPostCreated(postID int64, slug string, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "post.created",
		Entity:    "post",
		ID:        &postID,
		Slug:      &slug,
		Action:    "created",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerPostUpdated fires a webhook for post update
func TriggerPostUpdated(postID int64, slug string, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "post.updated",
		Entity:    "post",
		ID:        &postID,
		Slug:      &slug,
		Action:    "updated",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerPostDeleted fires a webhook for post deletion
func TriggerPostDeleted(postID int64, slug string) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "post.deleted",
		Entity:    "post",
		ID:        &postID,
		Slug:      &slug,
		Action:    "deleted",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      nil,
	}
	service.FireWebhook(payload)
}

// TriggerCategoryCreated fires a webhook for category creation
func TriggerCategoryCreated(categoryID int64, slug string, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "category.created",
		Entity:    "category",
		ID:        &categoryID,
		Slug:      &slug,
		Action:    "created",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerCategoryUpdated fires a webhook for category update
func TriggerCategoryUpdated(categoryID int64, slug string, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "category.updated",
		Entity:    "category",
		ID:        &categoryID,
		Slug:      &slug,
		Action:    "updated",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerCategoryDeleted fires a webhook for category deletion
func TriggerCategoryDeleted(categoryID int64, slug string) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "category.deleted",
		Entity:    "category",
		ID:        &categoryID,
		Slug:      &slug,
		Action:    "deleted",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      nil,
	}
	service.FireWebhook(payload)
}

// TriggerTagCreated fires a webhook for tag creation
func TriggerTagCreated(tagID int64, slug string, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "tag.created",
		Entity:    "tag",
		ID:        &tagID,
		Slug:      &slug,
		Action:    "created",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerTagUpdated fires a webhook for tag update
func TriggerTagUpdated(tagID int64, slug string, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "tag.updated",
		Entity:    "tag",
		ID:        &tagID,
		Slug:      &slug,
		Action:    "updated",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerTagDeleted fires a webhook for tag deletion
func TriggerTagDeleted(tagID int64, slug string) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "tag.deleted",
		Entity:    "tag",
		ID:        &tagID,
		Slug:      &slug,
		Action:    "deleted",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      nil,
	}
	service.FireWebhook(payload)
}

// TriggerAuthorCreated fires a webhook for author/user creation
func TriggerAuthorCreated(authorID int64, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "author.created",
		Entity:    "author",
		ID:        &authorID,
		Slug:      nil,
		Action:    "created",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerAuthorUpdated fires a webhook for author/user update
func TriggerAuthorUpdated(authorID int64, data map[string]interface{}) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "author.updated",
		Entity:    "author",
		ID:        &authorID,
		Slug:      nil,
		Action:    "updated",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      data,
	}
	service.FireWebhook(payload)
}

// TriggerAuthorDeleted fires a webhook for author/user deletion
func TriggerAuthorDeleted(authorID int64) {
	service := GetGlobalWebhookService()
	payload := models.WebhookPayload{
		Event:     "author.deleted",
		Entity:    "author",
		ID:        &authorID,
		Slug:      nil,
		Action:    "deleted",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Data:      nil,
	}
	service.FireWebhook(payload)
}
