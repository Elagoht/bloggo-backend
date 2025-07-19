package category

type CategoryHandler struct {
	service CategoryService
}

func NewCategoryHandler(service CategoryService) CategoryHandler {
	return CategoryHandler{
		service,
	}
}
