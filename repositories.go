package main

import (
	"log"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
	USERS_FIND_ALL_SQL                   = "select * from users"
	USERS_FIND_BY_ID_SQL                 = "select * from users where id = $1"
	USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL = "select * from users where email = lower($1) and password_digest = crypt($2, password_digest)"
	USERS_INSERT_SQL                     = "insert into users (name, email, password_digest) values ($1, lower($2), crypt($3, gen_salt('bf', 8))) returning *"
	SITES_FIND_ALL_SQL                   = "select * from sites"
	SITES_INSERT_SQL                     = "insert into sites (name) values ($1)"
)

type UserRepository struct {
	db *sql.DB
}

type VisitRepository struct {
	db *sql.DB
}

type SitesRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func NewVisitRepository(db *sql.DB) *VisitRepository {
	return &VisitRepository{
		db: db,
	}
}

func NewSitesRepository(db *sql.DB) *SitesRepository {
	return &SitesRepository{
		db: db,
	}
}

func (repo *UserRepository) FetchAll() []User {
	var users []User

	rows, _ := repo.db.Query(USERS_FIND_ALL_SQL)

	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.id, &user.Name, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, *user)
	}

	return users
}

func (repo *UserRepository) FindById(id int) *User {
	user := new(User)

	err := repo.db.QueryRow(USERS_FIND_BY_ID_SQL, id).
		Scan(&user.id, &user.Name, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			return nil
		} else {
			log.Fatal(err)
		}
	}

	return user
}

func (repo *UserRepository) FindByEmailAndPassword(email, password string) *User {
	user := new(User)

	err := repo.db.QueryRow(USERS_FIND_BY_EMAIL_AND_PASSWORD_SQL, email, password).
		Scan(&user.id, &user.Name, &user.email, &user.passwordDigest, &user.createdAt, &user.updatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// there were no rows, but otherwise no error occurred
			return nil
		} else {
			log.Fatal(err)
		}
	}

	return user
}

func (repo *SitesRepository) FindAll() []*Site {
	var sites []*Site

	rows, _ := repo.db.Query(SITES_FIND_ALL_SQL)

	for rows.Next() {
		site := new(Site)
		err := rows.Scan(&site.id, &site.name, &site.createdAt, &site.updatedAt)
		if err != nil {
			log.Fatal(err)
		}
		sites = append(sites, site)
	}

	return sites
}
