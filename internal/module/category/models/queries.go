package models

import "bloggo/internal/utils/slugify"

// -- Create Category Params -- //
type QueryParamsCategoryCreate struct {
	Name        string
	Slug        string
	Description string
}

func ToCreateCategoryParams(
	model *RequestCategoryCreate,
) *QueryParamsCategoryCreate {
	return &QueryParamsCategoryCreate{
		model.Name,
		slugify.Slugify(model.Name),
		model.Description,
	}
}

// -- Patch Category Params -- //
type QueryParamsCategoryUpdate struct {
	Name        *string
	Slug        *string
	Description *string
}

func ToUpdateCategoryParams(
	model *RequestCategoryUpdate,
) *QueryParamsCategoryUpdate {
	params := &QueryParamsCategoryUpdate{}

	if model.Name != "" {
		params.Name = &model.Name
		slug := slugify.Slugify(model.Name)
		params.Slug = &slug
	}

	if model.Description != "" {
		params.Description = &model.Description
	}

	return params
}
