package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/weatherglass/models"
)

const (
	USERS_FIND_ALL_SQL                   = "select * from users"
	USERS_FIND_BY_ID_SQL                 = "select * from users where id = $1"
	USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL = "select * from users where email = lower($1) and password_digest = crypt($2, password_digest)"
	USERS_INSERT_SQL                     = "insert into users (name, email, password_digest, email_confirmation_token) values ($1, lower($2), crypt($3, gen_salt('bf', 8)), $4) returning *"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (repo *UsersRepository) FetchAll() ([]models.User, error) {
	var users []models.User
	rows, err := repo.db.Query(USERS_FIND_ALL_SQL)

	for rows.Next() {
		user := new(models.User)
		err = user.FromRow(rows)
		users = append(users, *user)
	}

	return users, err
}

type UserResult struct {
	Ok  *models.User
	Err error
}

func (repo *UsersRepository) FindById(id int) <-chan UserResult {
	result := make(chan UserResult)

	go func() {
		user := new(models.User)
		row := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id)
		err := user.FromRow(row)
		if err != nil && err == sql.ErrNoRows {
			result <- UserResult{nil, err}
			return
		}

		result <- UserResult{user, nil}
	}()

	return result
}

func (repo *UsersRepository) FindByEmailAndPassword(email, password string) (*models.User, error) {
	user := new(models.User)
	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password)
	err := user.FromRow(row)
	if err != nil && err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (repo *UsersRepository) Create(user *models.User) (*models.User, error) {
	row := repo.db.QueryRow(USERS_INSERT_SQL, user.Name, user.Email,
		user.PasswordDigest, user.EmailConfirmationToken)
	err := user.FromRow(row)
	return user, err
}
