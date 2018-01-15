package main

import (
	"fmt"
	"time"
)

var currentId int

type Objects []interface{}

var people People
var objects Objects

// Create Some Default People At Init
func init() {
	RepoCreatePerson(Person{FirstName: "Mitchell", LastName: "Strong", IsUHN: true})
	RepoCreatePerson(Person{FirstName: "Richard", LastName: "de Borja", IsUHN: true})
}

func RepoFindPerson(id int) Person {
	for _, t := range people {
		if t.ID == id {
			return t
		}
	}
	// return empty Person if not found
	return Person{}
}

func RepoCreatePerson(t Person) Person {
	currentId++
	t.ID = currentId
	t.AddedOn = time.Now()
	people = append(people, t)
	return t
}

func RepoDestroyPerson(id int) error {
	for i, t := range people {
		if t.ID == id {
			people = append(people[:i], people[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Person with id of %d to delete", id)
}
