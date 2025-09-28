package models

type SearchResultType string

const (
	SearchResultTypeTag      SearchResultType = "tag"
	SearchResultTypeCategory SearchResultType = "category"
	SearchResultTypePost     SearchResultType = "post"
	SearchResultTypeUser     SearchResultType = "user"
)

type SearchResult struct {
	ID        int64            `json:"id"`
	Type      SearchResultType `json:"type"`
	Title     string           `json:"title"`
	Slug      *string          `json:"slug,omitempty"`
	AvatarURL *string          `json:"avatarUrl,omitempty"`
	CoverURL  *string          `json:"coverUrl,omitempty"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}