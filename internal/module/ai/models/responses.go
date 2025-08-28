package models

type ResponseGenerativeFill struct {
	Title             string `json:"title"`
	MetaDescription   string `json:"metaDescription"`
	Spot              string `json:"spot"`
	SuggestedCategory string `json:"suggestedCategory"`
}
