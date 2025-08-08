package models

import "bloggo/internal/utils/slugify"

// -- Create Tag Params -- //
type QueryParamsTagCreate struct {
	Name string
	Slug string
}

func ToCreateTagParams(
	model *RequestTagCreate,
) *QueryParamsTagCreate {
	return &QueryParamsTagCreate{
		model.Name,
		slugify.Slugify(model.Name),
	}
}

// -- Patch Tag Params -- //
type QueryParamsTagUpdate struct {
	Name *string
	Slug *string
}

func ToUpdateTagParams(
	model *RequestTagUpdate,
) *QueryParamsTagUpdate {
	params := &QueryParamsTagUpdate{}

	if model.Name != "" {
		params.Name = &model.Name
		slug := slugify.Slugify(model.Name)
		params.Slug = &slug
	}

	return params
}
