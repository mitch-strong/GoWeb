package main

import (
	"fmt"
	"time"
)

var currentID int

//Objects is an array of empty interfaces
type Objects []interface{}

var people People
var objects Objects

// Create Some Default People At Init
func init() {
	RepoCreatePerson(Person{FirstName: "Mitchell", LastName: "Strong", IsUHN: true})
	RepoCreatePerson(Person{FirstName: "Richard", LastName: "de Borja", IsUHN: true})
}

//RepoFindPerson Finds a person in the people array based on ID
func RepoFindPerson(id int) Person {
	for _, t := range people {
		if t.ID == id {
			return t
		}
	}
	// return empty Person if not found
	return Person{}
}

//RepoCreatePerson Adds a person to the people Array
func RepoCreatePerson(t Person) Person {
	currentID++
	t.ID = currentID
	t.AddedOn = time.Now()
	people = append(people, t)
	return t
}

//RepoDestroyPerson Removes a person from the people array based on ID
func RepoDestroyPerson(id int) error {
	for i, t := range people {
		if t.ID == id {
			people = append(people[:i], people[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Person with id of %d to delete", id)
}
