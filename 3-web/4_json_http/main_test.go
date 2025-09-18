package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestCreateUsers(t *testing.T) {
	t.Parallel()

	h := Handlers{
		users: []User{},
		mu:    &sync.Mutex{},
	}

	body := bytes.NewReader([]byte(`{"name": "Vasily", "password": "qwerty"}`))

	expectedUsers := []User{
		{
			ID:       3,
			Name:     "Vasily",
			Password: "qwerty",
		},
	}

	r := httptest.NewRequest("POST", "/users/", body)
	w := httptest.NewRecorder()

	h.HandleCreateUser(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	reflect.DeepEqual(h.users, expectedUsers)
}

var expectedJSON = `[{"id":1,"name":"Afanasiy"},{"id":2,"name":"Ka"}]`

func TestGetUsers(t *testing.T) {

	h := Handlers{
		users: []User{
			{
				ID:       1,
				Name:     "Afanasiy",
				Password: "1234",
			},
			{
				ID:       2,
				Name:     "Ka",
				Password: "jdjfaljhfljehfs;l3345354",
			},
		},
		mu: &sync.Mutex{},
	}

	t.Parallel()

	r := httptest.NewRequest("GET", "/users/", nil)
	w := httptest.NewRecorder()

	h.HandleListUsers(w, r)

	if w.Code != http.StatusOK {
		t.Error("status is not ok")
	}

	bytes, _ := io.ReadAll(w.Body)
	if strings.Trim(string(bytes), "\n") != expectedJSON {
		t.Errorf("expected: [%s], got: [%s]", expectedJSON, string(bytes))
	}
}
