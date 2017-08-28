package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Visit struct {
	Id        int `json:"id"`
	Location  int `json:"location"`
	User      int `json:"user"`
	VisitedAt int `json:"visited_at"`
	Mark      int `json:"mark"`
}

func newVisit(w http.ResponseWriter, r *http.Request) {
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

	var body Visit
	err = json.Unmarshal(data, &body)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	err = DB.NewVisit(&body)
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func getVisit(w http.ResponseWriter, r *http.Request, id int) {
	visit, err := DB.GetVisit(id)
	if err == NotFound {
		Log.Warnf("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get id %s. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	result, err := json.Marshal(visit)
	if err != nil {
		Log.Errorf("Cannot marshal visit %#v. Reason %s", visit, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal visit"))
		return
	}

	writeAnswer(w, http.StatusOK, string(result))
}

func updateVisit(w http.ResponseWriter, r *http.Request, id int) {
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

	visit, err := DB.GetVisit(id)
	if err == NotFound {
		Log.Infof("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get visit with id %d. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot get visit"))
		return
	}

	oldUser := visit.User
	oldLocation := visit.Location

	err = json.Unmarshal(data, &visit)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	//UPDATE USER
	if visit.User != oldUser {
		usrOld, err := DB.GetUser(oldUser)
		if err != nil {
			Log.Errorf("Cannot get user %d. Reason %s", oldUser, err)
			writeAnswer(w, http.StatusBadRequest, "")
			return
		}
		usrOld.Visits.Remove(visit.Id)

		usr, err := DB.GetUser(visit.User)
		if err != nil {
			Log.Errorf("Cannot get user %d. Reason %s", visit.User, err)
			writeAnswer(w, http.StatusBadRequest, "")
			return
		}
		usr.Visits.Add(visit.Id)
	}

	//UPDATE LOCATION
	if visit.Location != oldLocation {
		locOld, err := DB.GetLocation(oldLocation)
		if err != nil {
			Log.Errorf("Cannot get location %d. Reason %s", oldLocation, err)
			writeAnswer(w, http.StatusBadRequest, "")
			return
		}
		locOld.Visits.Remove(visit.Id)

		loc, err := DB.GetLocation(visit.Location)
		if err != nil {
			Log.Errorf("Cannot get location %d. Reason %s", visit.Location, err)
			writeAnswer(w, http.StatusBadRequest, "")
			return
		}
		loc.Visits.Add(visit.Id)
	}

	err = DB.UpdateVisit(visit, id)
	if err != nil {
		Log.Errorf("Cannot update visit. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot update visit"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func processVisit(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		Log.Infof("Cannot parse id %s. Reason %s", path[2], err)
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if r.Method == "GET" {
		getVisit(w, r, id)
	} else {
		updateVisit(w, r, id)
	}
}
