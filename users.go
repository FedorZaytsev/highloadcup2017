package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Birthdate int    `json:"birth_date"`
	Visits    Array  `json:"-"`
}

type UserVisits struct {
	VisitedAt int    `json:"visited_at"`
	Mark      int    `json:"mark"`
	Place     string `json:"place"`
}

type UserVisitsSorter struct {
	Data []UserVisits
}

func (s UserVisitsSorter) Len() int {
	return len(s.Data)
}

func (s UserVisitsSorter) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func (s UserVisitsSorter) Less(i, j int) bool {
	return s.Data[i].VisitedAt < s.Data[j].VisitedAt
}

func newUser(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Warnf("Cannot read body from request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot read body from request"))
		return
	}

	if strings.Index(string(data), "\": null") != -1 {
		Log.Infof("null param")
		writeAnswer(w, http.StatusBadRequest, "")
		return
	}

	body := User{
		Visits: NewArray(),
	}
	Log.Infof("Visits %p", body.Visits)
	err = json.Unmarshal(data, &body)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}
	Log.Infof("Visits after %p", body.Visits)

	err = DB.NewUser(&body)
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func getUserVisits(w http.ResponseWriter, r *http.Request, id int) {
	filters := r.URL.Query()
	Log.Infof("Getting user %s visits with filters %#v", id, filters)

	vals, err := DB.GetVisitsFilter(id, filters)
	switch err {
	case NotFound:
		Log.Infof("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	case CannotParse:
		Log.Infof("Cannot parse")
		writeAnswer(w, http.StatusBadRequest, "")
		return
	case nil:
		break
	default:
		Log.Errorf("Error while getting user visits. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Error while getting user visits"))
		return
	}

	result, err := json.Marshal(struct {
		Visits []UserVisits `json:"visits"`
	}{
		Visits: vals,
	})
	if err != nil {
		Log.Errorf("Cannot marshal answer. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal answer"))
		return
	}

	writeAnswer(w, http.StatusOK, string(result))
}

func getUser(w http.ResponseWriter, r *http.Request, id int) {
	result, err := DB.GetUser(id)
	if err == NotFound {
		Log.Infof("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get id %d. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	answer, err := json.Marshal(result)
	if err != nil {
		Log.Errorf("Cannot marshal user %d. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal user"))
		return
	}

	writeAnswer(w, http.StatusOK, string(answer))
}

func updateUser(w http.ResponseWriter, r *http.Request, id int) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Warnf("Cannot read body from request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot read body from request"))
		return
	}

	if strings.Index(string(data), "\": null") != -1 {
		Log.Infof("null param")
		writeAnswer(w, http.StatusBadRequest, "")
		return
	}

	user, err := DB.GetUser(id)
	if err == NotFound {
		Log.Infof("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get user with id %d. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot get user"))
		return
	}

	err = json.Unmarshal(data, &user)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	err = DB.UpdateUser(user, id)
	if err != nil {
		Log.Errorf("Cannot update user. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot update user"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func processUser(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		Log.Infof("Cannot parse id %s. Reason %s", path[2], err)
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if path[len(path)-1] == "visits" {
		getUserVisits(w, r, id)
		return
	}
	if r.Method == "GET" {
		getUser(w, r, id)
	} else {
		updateUser(w, r, id)
	}
}
