package models

import "bloggo/internal/utils/slugify"

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
