package user

import (
	"bloggo/internal/module/user/models"
	"bloggo/internal/utils/filter"
	"bloggo/internal/utils/handlers"
	"bloggo/internal/utils/pagination"
	"database/sql"
	"fmt"
)

type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(database *sql.DB) UserRepository {
	return UserRepository{
		database,
	}
}

func (repository *UserRepository) GetUsers(
	paginate *pagination.PaginationOptions,
	search *filter.SearchOptions,
) ([]models.ResponseUserCard, error) {
	// Handle pagination and order params
	orderByClause, limitClause, offsetClause, args := paginate.BuildPaginationClauses()

	// Handle search by name
	searchClause, searchArgs := filter.BuildSearchClause(search, []string{"users.name"})

	// Merge them and generate query
	query, allArgs := handlers.BuildModifiedSQL(
		QueryUserGetUserCards,
		[]string{searchClause, orderByClause, limitClause, offsetClause},
		[][]any{searchArgs, args},
	)

	// Run query
	rows, err := repository.database.Query(query, allArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.ResponseUserCard{}
	for rows.Next() {
		var user models.ResponseUserCard
		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Avatar,
			&user.RoleId,
			&user.RoleName,
			&user.WrittenPostCount,
			&user.PublishedPostCount,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (repository *UserRepository) GetUsersCount(
	search *filter.SearchOptions,
) (int64, error) {
	// Handle search by name and email
	searchClause, searchArgs := filter.BuildSearchClause(search, []string{"users.name", "users.email"})

	// Generate query
	query, allArgs := handlers.BuildModifiedSQL(
		QueryUserCount,
		[]string{searchClause},
		[][]any{searchArgs},
	)

	// Execute query
	var count int64
	err := repository.database.QueryRow(query, allArgs...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *UserRepository) GetUserById(
	id int64,
) (*models.ResponseUserDetails, error) {
	row := repository.database.QueryRow(QueryUserGetById, id)

	var user models.ResponseUserDetails
	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Avatar,
		&user.CreatedAt,
		&user.LastLogin,
		&user.RoleId,
		&user.RoleName,
		&user.WrittenPostCount,
		&user.PublishedPostCount,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) UserCreate(
	model *models.UserCreateParams,
) (int64, error) {
	statement, err := repository.database.Prepare(QueryUserCreate)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(
		&model.Name,
		&model.Email,
		&model.Avatar,
		&model.PassphraseHash,
		&model.RoleId,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repository *UserRepository) UpdateAvatarById(
	userId int64,
	fileName string,
) error {
	statement, err := repository.database.Prepare(QueryUserUpdateAvatarById)
	if err != nil {
		return err
	}

	_, err = statement.Exec(fileName, userId)
	if err != nil {
		return err
	}

	return nil
}

func (repository *UserRepository) UpdateUserById(
	userId int64,
	model *models.RequestUserUpdate,
) error {
	_, err := repository.database.Exec(
		QueryUserUpdateById,
		model.Name,
		model.Email,
		userId,
	)
	return err
}

func (repository *UserRepository) AssignRole(
	userId int64,
	roleId int64,
) error {
	_, err := repository.database.Exec(QueryUserAssignRole, roleId, userId)
	return err
}

func (repository *UserRepository) DeleteUser(userId int64) error {
	_, err := repository.database.Exec(QueryUserDelete, userId)
	return err
}

func (repository *UserRepository) UpdateLastLogin(userId int64) error {
	_, err := repository.database.Exec(QueryUserUpdateLastLogin, userId)
	return err
}

func (repository *UserRepository) UpdatePasswordById(
	userId int64,
	hashedPassword string,
) error {
	_, err := repository.database.Exec(QueryUserUpdatePasswordById, hashedPassword, userId)
	return err
}
