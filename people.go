package main

import "time"

//Person Object with JSON definitions
type Person struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first"`
	LastName  string    `json:"last"`
	IsUHN     bool      `json:"UHN"`
	AddedOn   time.Time `json:"AddedOn"`
}

//People is a list of Persons
type People []Person
