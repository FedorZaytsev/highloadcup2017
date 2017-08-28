package main

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type Location struct {
	Id       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`
	Visits   Array  `json:"-"`
}

func newLocation(w http.ResponseWriter, r *http.Request) {
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

	body := Location{
		Visits: NewArray(),
	}
	err = json.Unmarshal(data, &body)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	err = DB.NewLocation(&body)
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func getLocationAvg(w http.ResponseWriter, r *http.Request, id int) {
	filters := r.URL.Query()
	Log.Infof("Getting avg for location %d with filters %#v", id, filters)

	val, err := DB.GetAverage(id, filters)
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
		Avg float64 `json:"avg"`
	}{
		Avg: toFixed(float64(val), 5),
	})
	if err != nil {
		Log.Errorf("Cannot marshal answer. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal answer"))
		return
	}

	writeAnswer(w, http.StatusOK, string(result))
}

func getLocation(w http.ResponseWriter, r *http.Request, id int) {
	location, err := DB.GetLocation(id)
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

	result, err := json.Marshal(location)
	if err != nil {
		Log.Errorf("Cannot marshal location %#v. Reason %s", location, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal location"))
		return
	}

	writeAnswer(w, http.StatusOK, string(result))
}

func updateLocation(w http.ResponseWriter, r *http.Request, id int) {
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

	loc, err := DB.GetLocation(id)
	if err == NotFound {
		Log.Infof("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get location with id %d. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot get location"))
		return
	}

	err = json.Unmarshal(data, &loc)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	err = DB.UpdateLocation(loc, id)
	if err != nil {
		Log.Errorf("Cannot update location. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot update location"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func processLocation(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	path := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		Log.Infof("Cannot parse id %s. Reason %s", path[2], err)
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if path[len(path)-1] == "avg" {
		getLocationAvg(w, r, id)
		return
	}
	if r.Method == "GET" {
		getLocation(w, r, id)
	} else {
		updateLocation(w, r, id)
	}
}
