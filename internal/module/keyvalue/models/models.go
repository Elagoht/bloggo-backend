package models

type KeyValue struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type RequestKeyValueUpsert struct {
	Key   string `json:"key" validate:"required,max=255"`
	Value string `json:"value" validate:"required"`
}

type RequestKeyValueBulkUpsert struct {
	Items []RequestKeyValueUpsert `json:"items" validate:"required,dive"`
}
