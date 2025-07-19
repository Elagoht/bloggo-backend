package category

type CategoryService struct {
	repository CategoryRepository
}

func NewCategoryService(repository CategoryRepository) CategoryService {
	return CategoryService{
		repository,
	}
}

func (service CategoryService) ListCategories() ([]Category, error) {
	categories, errorFromDatabase := service.repository.GetAllCategories()
	if errorFromDatabase != nil {
		return nil, errorFromDatabase
	}
	return categories, nil
}

func (service CategoryService) CreateCategory(
	categoryInput CreateCategoryRequest,
) (Category, error) {
	categoryToSave := Category{
		Title:       categoryInput.Title,
		Description: categoryInput.Description,
	}

	savedCategory, errorFromDatabase := service.repository.CreateCategory(categoryToSave)
	if errorFromDatabase != nil {
		return Category{}, errorFromDatabase
	}

	return savedCategory, nil
}

func (service CategoryService) UpdateCategory(
	categoryInput UpdateCategoryRequest,
) error {
	categoryToUpdate := Category{
		ID:          categoryInput.ID,
		Title:       categoryInput.Title,
		Description: categoryInput.Description,
	}

	errorFromDatabase := service.repository.UpdateCategory(categoryToUpdate)
	if errorFromDatabase != nil {
		return errorFromDatabase
	}

	return nil
}

func (service CategoryService) DeleteCategory(
	categoryIdentifier string,
) error {
	errorFromDatabase := service.repository.DeleteCategory(categoryIdentifier)
	if errorFromDatabase != nil {
		return errorFromDatabase
	}

	return nil
}
