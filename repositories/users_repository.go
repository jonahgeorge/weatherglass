package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/weatherglass/models"
)

const (
	USERS_FIND_ALL_SQL                   = "select * from users"
	USERS_FIND_BY_ID_SQL                 = "select * from users where id = $1"
	USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL = "select * from users where email = lower($1) and password_digest = crypt($2, password_digest)"
	USERS_INSERT_SQL                     = "insert into users (name, email, password_digest) values ($1, lower($2), crypt($3, gen_salt('bf', 8))) returning *"
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

func (repo *UsersRepository) FindById(id int) (*models.User, error) {
	user := new(models.User)

	row := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id)
	err := user.FromRow(row)

	if err != nil && err == sql.ErrNoRows {
		// there were no rows, but otherwise no error occurred
		return nil, nil
	}

	return user, err
}

func (repo *UsersRepository) FindByEmailAndPassword(email, password string) (*models.User, error) {
	user := new(models.User)

	row := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password)
	err := user.FromRow(row)

	if err != nil && err == sql.ErrNoRows {
		// there were no rows, but otherwise no error occurred
		return nil, nil
	}

	return user, err
}
