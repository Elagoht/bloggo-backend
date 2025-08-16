package models

type RequestCreateRemovalRequest struct {
	PostVersionId int64  `json:"postVersionId" validate:"required"`
	Note          string `json:"note" validate:"max=500"`
}

type RequestDecideRemovalRequest struct {
	Note string `json:"note" validate:"max=500"`
}