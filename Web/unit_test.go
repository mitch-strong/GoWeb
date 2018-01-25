package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"
)

//Currently all authentication is skipped during unit testing
func TestMain(m *testing.M) {
	goTest = true
	router := NewRouter()
	go http.ListenAndServe(":"+localport, router)
	os.Exit(m.Run())
}

func TestUsersAdded(t *testing.T) {
	if len(people) != 2 {
		t.Fail()
		return
	}
	for i := range people {
		if people[i].FirstName == "" || people[i].LastName == "" || people[i].FirstName == "Homer" {
			t.Fail()
			return
		}

		if !people[i].IsUHN {
			t.Fail()
			return
		}
		var emptyTime time.Time
		if people[i].AddedOn == emptyTime {
			t.Fail()
			return
		}
	}
}

//TestAddNewUserPOST Tests if a new user can be added via API POST
func TestAddNewUserPOST(t *testing.T) {
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
		t.Fail()
		return
	}
	body := bytes.NewReader(payloadBytes)
	test := server + "/people"
	_ = test
	req, err := http.NewRequest("POST", server+"/people", body)
	if err != nil {
		t.Fail()
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Log(err)
		t.Fail()
		return
	}
	defer resp.Body.Close()

	if len(people) != 3 {
		t.Fail()
		return
	}
	if people[2].FirstName == "" || people[2].LastName == "" {
		t.Fail()
		return
	}

	if people[2].IsUHN {
		t.Fail()
		return
	}
	var emptyTime time.Time
	if people[2].AddedOn == emptyTime {
		t.Fail()
		return
	}
}

//AddJSONsPOST adds two JSON objects to the JSON array via a POST to /JSON
func TestAddJSONsPOST(t *testing.T) {
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
		t.Fail()
		return
	}

	if err2 != nil {
		t.Fail()
		return
	}

	body := bytes.NewReader(payloadBytes)
	body2 := bytes.NewReader(payloadBytes2)

	req, err := http.NewRequest("POST", server+"/JSON", body)
	req2, err2 := http.NewRequest("POST", server+"/JSON", body2)

	if err != nil {
		t.Fail()
		return
	}
	if err2 != nil {
		t.Fail()
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	resp2, err2 := http.DefaultClient.Do(req2)

	if err != nil {
		t.Fail()
		return
	}
	if err2 != nil {
		t.Fail()
		return
	}
	defer resp.Body.Close()
	defer resp2.Body.Close()

	if len(objects) != 2 {
		t.Fail()
		return
	}
}

func TestUrls(t *testing.T) {
	result := true
	result = result && CheckPageResponse(server+"/")
	result = result && CheckPageResponse(server+"/people")
	result = result && CheckPageResponse(server+"/peopleJSON")
	result = result && CheckPageResponse(server+"/JSON")
	result = result && CheckNoPageResponse(server+"/x")

	if result != true {
		t.Fail()
	}
}

//DeleteUser tests if the new user can be deleted (Dependency on previous)
func TestDeleteUser(t *testing.T) {
	err := RepoDestroyPerson(3)
	if err != nil {
		t.Fail()
		return
	}
	t.Run("Check if Deleted", TestUsersAdded)
}

//CheckPageResponse checks if a page that should respond is found correctly
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

//CheckNoPageResponse checks if a page that does not exist responds with a 404 Error
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
