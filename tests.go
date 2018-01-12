package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var pass, fail int

func runUnitTests() {
	pass, fail = 0, 0
	log.Printf("////////////////////////////////////////////////////////")
	log.Printf("Beginning to run tests")
	log.Printf("Testing: Default Users Added - %v", PassFail(testUsersAdded()))
	log.Printf("Testing: Adding User By JSON POST - %v", PassFail(AddNewUserPOST()))
	log.Printf("Testing: Deleting New User - %v", PassFail(DeleteUser(3)))
	log.Printf("Testing: Check if '/' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/")))
	log.Printf("Testing: Check if '/person' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/people")))
	log.Printf("Testing: Check if '/JSONperson' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/JSONpeople")))
	log.Printf("PASS - %v    FAIL - %v", pass, fail)
	log.Printf("////////////////////////////////////////////////////////")
}

//Tests if Default Users Are Added Correctly
func testUsersAdded() bool {
	if len(people) != 2 {
		return false
	}
	for i := range people {
		if people[i].FirstName == "" || people[i].LastName == "" || people[i].FirstName == "Homer" {
			return false
		}

		if !people[i].IsUHN {
			return false
		}
		var emptyTime time.Time
		if people[i].AddedOn == emptyTime {
			return false
		}
	}
	return true
}

//Tests if a new user can be added via API POST
func AddNewUserPOST() bool {
	type Payload struct {
		First string `json:"first"`
		Last  string `json:"last"`
		UHN   bool   `json:"UHN"`
	}

	data := Payload{
		First: "Homer",
		Last:  "Simpson",
		UHN:   false,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return false
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:8080/people", body)
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// if err != nil {
	// 	log.Printf("%v", err)
	// 	return false
	// }

	if len(people) != 3 {
		return false
	}
	if people[2].FirstName == "" || people[2].LastName == "" {
		return false
	}

	if people[2].IsUHN {
		return false
	}
	var emptyTime time.Time
	if people[2].AddedOn == emptyTime {
		return false
	}

	return true
}

//Tests if the new user can be deleted (Dependency on previous)
func DeleteUser(id int) bool {
	err := RepoDestroyPerson(3)
	if err != nil {
		return false
	}
	if !testUsersAdded() {
		return false
	}
	return true

}

func CheckPageResponse(url string) bool {
	response, err := http.Get(url)
	if err != nil {
		return false
	}
	if response == nil {
		return false
	}
	return true
}

//Returns Pass or Fail
func PassFail(result bool) string {
	if result {
		pass++
		return "Pass"
	}
	fail++
	return "Fail"
}
