package repositories

import "database/sql"

type VisitRepository struct {
	db *sql.DB
}

func NewVisitRepository(db *sql.DB) *VisitRepository {
	return &VisitRepository{db: db}
}
