package main

import "time"

type User struct {
	id int
	Name string
	email string
	passwordDigest string
	createdAt time.Time
	updatedAt time.Time
}

type Site struct {
	id int
	name string
	createdAt time.Time
	updatedAt time.Time
}