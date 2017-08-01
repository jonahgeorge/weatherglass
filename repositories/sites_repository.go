package repositories

import (
	"database/sql"

	"github.com/jonahgeorge/weatherglass/models"
)

const (
	SITES_FIND_BY_ID_SQL      = "select * from sites where id = $1"
	SITES_FIND_BY_USER_ID_SQL = "select * from sites where user_id = $1"
	SITES_INSERT_SQL          = "insert into sites (user_id, name) values ($1, $2)"
	SITES_UPDATE_SQL          = "update sites set name $2 where id = $1"
	SITES_DELETE_SQL          = "delete from sites where id = $1"
)

type SitesRepository struct {
	db *sql.DB
}

func NewSitesRepository(db *sql.DB) *SitesRepository {
	return &SitesRepository{db: db}
}

func (repo *SitesRepository) FindByUserId(userId int) ([]models.Site, error) {
	var sites []models.Site
	rows, err := repo.db.Query(SITES_FIND_BY_USER_ID_SQL, userId)

	for rows.Next() {
		site := new(models.Site)
		err = site.FromRow(rows)
		sites = append(sites, *site)
	}

	return sites, err
}

func (repo *SitesRepository) FindById(id int) (*models.Site, error) {
	site := new(models.Site)
	row := repo.db.QueryRow(SITES_FIND_BY_ID_SQL, id)
	err := site.FromRow(row)
	return site, err
}

func (repo *SitesRepository) Create(userId int, name string) (*models.Site, error) {
	site := new(models.Site)
	row := repo.db.QueryRow(SITES_INSERT_SQL, userId, name)
	err := site.FromRow(row)
	return site, err
}

func (repo *SitesRepository) Update(id int, name string) (*models.Site, error) {
	site := new(models.Site)
	row := repo.db.QueryRow(SITES_UPDATE_SQL, id, name)
	err := site.FromRow(row)
	return site, err
}

func (repo *SitesRepository) Delete(id int) (sql.Result, error) {
	return repo.db.Exec(SITES_DELETE_SQL, id)
}
