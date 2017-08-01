package models

import "time"

type Site struct {
	Id        int
	UserId    int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Site) CanCreate(u *User) bool {
	return s.UserId == u.Id
}

func (s *Site) CanView(u *User) bool {
	return s.UserId == u.Id
}

func (s *Site) CanUpdate(u *User) bool {
	return s.UserId == u.Id
}

func (s *Site) CanDelete(u *User) bool {
	return s.UserId == u.Id
}

func (s *Site) FromRow(row Scannable) error {
	return row.Scan(&s.Id, &s.UserId, &s.Name, &s.CreatedAt, &s.UpdatedAt)
}
