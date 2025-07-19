package category

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) CategoryService {
	return CategoryService{
		repository,
	}
}
