package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type Location struct {
	Id       int    `json:"id"`
	Place    string `json:"place"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Distance int    `json:"distance"`
}

func newLocation(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Warnf("Cannot read body from request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot read body from request"))
		return
	}

	var body Location
	err = json.Unmarshal(data, &body)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	err = RedisClientLoc.Set(strconv.Itoa(body.Id), data, 0).Err()
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func getLocationAvg(w http.ResponseWriter, r *http.Request, id string) {
}

func getLocation(w http.ResponseWriter, r *http.Request, id string) {
	Log.Infof("Getting location with id %s", id)

	result, err := RedisClientLoc.Get(id).Result()
	if err == redis.Nil {
		Log.Warnf("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get id %s. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	writeAnswer(w, http.StatusOK, result)
}

func updateLocation(w http.ResponseWriter, r *http.Request, id string) {
	Log.Infof("Updaing location with id %s", id)

	result, err := RedisClientLoc.Get(id).Result()
	if err == redis.Nil {
		Log.Warnf("Not found")
		writeAnswer(w, http.StatusNotFound, "")
		return
	}
	if err != nil {
		Log.Errorf("Cannot get id %s. Reason %s", id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	var loc Location
	err = json.Unmarshal([]byte(result), &loc)
	if err != nil {
		Log.Warnf("Cannot parse JSON from redis. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot parse JSON from redis"))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Warnf("Cannot read body from request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot read body from request"))
		return
	}

	var newLoc Location
	err = json.Unmarshal(data, &newLoc)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	if newLoc.City != "" {
		loc.City = newLoc.City
	}
	if newLoc.Country != "" {
		loc.Country = newLoc.Country
	}
	if newLoc.Place != "" {
		loc.Place = newLoc.Place
	}
	if newLoc.Distance != 0 {
		loc.Distance = newLoc.Distance
	}

	resultRedis, err := json.Marshal(loc)
	if err != nil {
		Log.Errorf("Cannot marshal json. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal json"))
		return
	}

	err = RedisClientLoc.Set(id, resultRedis, 0).Err()
	if err != nil {
		Log.Errorf("Cannot set redis value. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set redis value"))
		return
	}

	writeAnswer(w, http.StatusOK, "")
}

func processLocation(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[2]
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
