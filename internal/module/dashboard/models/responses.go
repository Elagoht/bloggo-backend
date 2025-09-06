package models

type ResponseDashboardStats struct {
	// Pending content section
	PendingVersions   []PendingVersion    `json:"pendingVersions"`
	RecentActivity    []RecentActivity    `json:"recentActivity"`
	PublishingRate    PublishingRate      `json:"publishingRate"`
	AuthorPerformance []AuthorPerformance `json:"authorPerformance"`

	// Content management
	DraftCount   DraftCount   `json:"draftCount"`
	PopularTags  []PopularTag `json:"popularTags"`
	StorageUsage StorageUsage `json:"storageUsage"`
}

type PendingVersion struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	AuthorId   int64  `json:"authorId"`
	AuthorName string `json:"authorName"`
	CreatedAt  string `json:"createdAt"`
}

type RecentActivity struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

type PublishingRate struct {
	ThisWeek  int `json:"thisWeek"`
	ThisMonth int `json:"thisMonth"`
}

type AuthorPerformance struct {
	AuthorId   int64  `json:"authorId"`
	AuthorName string `json:"authorName"`
	PostCount  int    `json:"postCount"`
}

type DraftCount struct {
	TotalDrafts    int              `json:"totalDrafts"`
	DraftsByAuthor []DraftsByAuthor `json:"draftsByAuthor"`
}

type DraftsByAuthor struct {
	AuthorId   int64  `json:"authorId"`
	AuthorName string `json:"authorName"`
	DraftCount int    `json:"draftCount"`
}

type PopularTag struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Usage int    `json:"usage"`
}

type StorageUsage struct {
	TotalSizeBytes int64   `json:"totalSizeBytes"`
	TotalSizeMB    float64 `json:"totalSizeMB"`
	FileCount      int     `json:"fileCount"`
}
