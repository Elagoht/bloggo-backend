package models

type AuditLogEntry struct {
	UserID     *int64                 `json:"userId,omitempty"`
	EntityType string                 `json:"entityType"`
	EntityID   int64                  `json:"entityId"`
	Action     string                 `json:"action"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Predefined action constants
const (
	// Basic CRUD actions
	ActionCreated = "created"
	ActionUpdated = "updated"
	ActionDeleted = "deleted"

	// Authentication actions
	ActionLogin  = "login"
	ActionLogout = "logout"

	// Workflow actions
	ActionSubmitted   = "submitted"
	ActionApproved    = "approved"
	ActionRejected    = "rejected"
	ActionPublished   = "published"
	ActionUnpublished = "unpublished"

	// Special actions
	ActionAssigned         = "assigned"
	ActionRemoved          = "removed"
	ActionDuplicatedFrom   = "duplicated_from"
	ActionReplacedPublished = "replaced_published"
	ActionRequested        = "requested"
	ActionAdded           = "added"
	ActionDenied          = "denied"

	// Legacy constants for backward compatibility (deprecated)
	ActionUserCreated = ActionCreated
	ActionUserUpdated = ActionUpdated
	ActionUserDeleted = ActionDeleted
	ActionUserLogin   = ActionLogin
	ActionUserLogout  = ActionLogout

	ActionPostCreated = ActionCreated
	ActionPostDeleted = ActionDeleted

	ActionVersionCreated           = ActionCreated
	ActionVersionUpdated           = ActionUpdated
	ActionVersionDeleted           = ActionDeleted
	ActionVersionSubmitted         = ActionSubmitted
	ActionVersionApproved          = ActionApproved
	ActionVersionRejected          = ActionRejected
	ActionVersionPublished         = ActionPublished
	ActionVersionUnpublished       = ActionUnpublished
	ActionVersionDuplicatedFrom    = ActionDuplicatedFrom
	ActionVersionReplacedPublished = ActionReplacedPublished

	ActionCategoryCreated = ActionCreated
	ActionCategoryUpdated = ActionUpdated
	ActionCategoryDeleted = ActionDeleted

	ActionTagCreated   = ActionCreated
	ActionTagUpdated   = ActionUpdated
	ActionTagDeleted   = ActionDeleted
	ActionTagsAssigned = ActionAssigned
	ActionTagsRemoved  = ActionRemoved

	ActionRoleCreated       = ActionCreated
	ActionRoleUpdated       = ActionUpdated
	ActionRoleDeleted       = ActionDeleted
	ActionPermissionAdded   = ActionAdded
	ActionPermissionRemoved = ActionRemoved

	ActionRemovalRequested = ActionRequested
	ActionRemovalApproved  = ActionApproved
	ActionRemovalDenied    = ActionDenied
)

// Entity type constants
const (
	EntityAuth           = "auth"
	EntityUser           = "user"
	EntityPost           = "post"
	EntityPostVersion    = "post_version"
	EntityCategory       = "category"
	EntityTag            = "tag"
	EntityRole           = "role"
	EntityPermission     = "permission"
	EntityRemovalRequest = "removal_request"
	EntityKeyValue       = "keyvalue"
)