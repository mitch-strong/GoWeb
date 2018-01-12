package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	const tpl = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>GoWeb</title>
	</head>
	<body>
		<h1>Welcome to GoWeb!</h1>
		<p>Please Note the Following Pages to Access</p>
		<ul>
			<li>GET --> "/" | Home Page</li>
			<li>GET --> "/JSONpeople" | Shows list of people as JSON strings</li>
			<li>GET --> "/people" | Shows list of people readable table</li>
			<li>POST --> "/people" | Creates a new person and adds them to the list
			<ul><li>"first" - string</li><li>"last" - string</li><li>"UHN" - bool</li></ul></li>
		</ul>
	</body>
	</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, tpl)
	return
}

func PersonList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
	<h1>Meet our current GoWeb Users!</h1>
	<table>
	<tr><th>First Name</th><th>Last Name</th><th>UHN Employee</th><th>Member Since</th></tr>
		{{range .PeopleList}}<tr><td>{{ .FirstName }}</td><td>{{ .LastName }}</td><td>{{ .IsUHN }}</td><td>{{ .AddedOn }}</td></tr>{{else}}<div><strong>No People</strong></div>{{end}}
	</table>
	</body>
</html>`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(tpl)
	check(err)

	data := struct {
		Title      string
		PeopleList []Person
	}{
		Title:      "List of People (Static)",
		PeopleList: people,
	}
	err = t.Execute(w, data)
	//err = t.Execute(os.Stdout, data)
	check(err)

	//fmt.Fprint(w, t)
	return
}

func PersonListJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(people); err != nil {
		panic(err)
	}
}

func PersonCreate(w http.ResponseWriter, r *http.Request) {
	var person Person
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &person); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreatePerson(person)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
}
