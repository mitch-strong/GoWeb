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
	x := PassFail(testUsersAdded())
	//log.Printf("Testing: Default Users Added - %v", x)
	test := Test{
		Name:   "Default Users Added",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(AddNewUserPOST())
	//log.Printf("Testing: Adding User By JSON POST - %v", PassFail(AddNewUserPOST()))
	test = Test{
		Name:   "Adding User By JSON POST",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(DeleteUser(3))
	//log.Printf("Testing: Deleting New User - %v", PassFail(DeleteUser(3)))
	test = Test{
		Name:   "Deleting New User",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(CheckPageResponse("http://localhost:8080/"))
	//log.Printf("Testing: Check if '/' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/")))
	test = Test{
		Name:   "Check if '/' Page Works",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(CheckPageResponse("http://localhost:8080/people"))
	//	log.Printf("Testing: Check if '/person' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/people")))
	test = Test{
		Name:   "Check if '/person' Page Works",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(CheckPageResponse("http://localhost:8080/peopleJSON"))
	//	log.Printf("Testing: Check if '/personJSON' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/peopleJSON")))
	test = Test{
		Name:   "Check if '/personJSON' Page Works",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(CheckPageResponse("http://localhost:8080/JSON"))
	//	log.Printf("Testing: Check if '/JSON' Page Works - %v", PassFail(CheckPageResponse("http://localhost:8080/JSON")))
	test = Test{
		Name:   "Check if '/JSON' Page Works",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(CheckNoPageResponse("http://localhost:8080/x"))
	//	log.Printf("Testing: Check for a Page that doesn't Exist - %v", PassFail(CheckNoPageResponse("http://localhost:8080/x")))
	test = Test{
		Name:   "Check for a Page that doesn't Exist",
		Status: x,
	}
	tests = append(tests, test)
	x = PassFail(AddJSONsPOST())
	//	log.Printf("Testing: Adding Two Random JSON Objects POST - %v", PassFail(AddJSONsPOST()))
	test = Test{
		Name:   "Adding Two Random JSON Objects POST",
		Status: x,
	}
	tests = append(tests, test)

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

func AddJSONsPOST() bool {
	type Payload struct {
		One   string `json:"Item1"`
		Two   string `json:"Item2"`
		Three bool   `json:"Item3"`
	}

	type Payload2 struct {
		One   string `json:"Item1"`
		Two   string `json:"Item2"`
		Three bool   `json:"Item3"`
		Four  int    `json:"Item4"`
		Five  string `json:"Item5"`
	}

	data := Payload{
		One:   "Homer",
		Two:   "Simpson",
		Three: false,
	}

	data2 := Payload2{
		One:   "This is a test",
		Two:   "Longstringtestwithnospaces",
		Three: false,
		Four:  7839,
		Five:  "15",
	}

	payloadBytes, err := json.Marshal(data)
	payloadBytes2, err2 := json.Marshal(data2)

	if err != nil {
		return false
	}

	if err2 != nil {
		return false
	}

	body := bytes.NewReader(payloadBytes)
	body2 := bytes.NewReader(payloadBytes2)

	req, err := http.NewRequest("POST", "http://localhost:8080/JSON", body)
	req2, err2 := http.NewRequest("POST", "http://localhost:8080/JSON", body2)

	if err != nil {
		return false
	}
	if err2 != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	resp2, err2 := http.DefaultClient.Do(req2)

	if err != nil {
		return false
	}
	if err2 != nil {
		return false
	}
	defer resp.Body.Close()
	defer resp2.Body.Close()

	if len(objects) != 2 {
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
	if response.Status == "404 Not Found" {
		return false
	}
	return true
}

func CheckNoPageResponse(url string) bool {
	response, err := http.Get(url)
	if err != nil {
		return true
	}
	if response == nil {
		return true
	}
	if response.Status == "404 Not Found" {
		return true
	}
	return false
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
