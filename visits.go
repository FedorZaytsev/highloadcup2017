package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
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

	var body Visit
	err = json.Unmarshal(data, &body)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	id := strconv.Itoa(body.Id)
	err = RedisClientVstCon.SAdd(strconv.Itoa(body.User), id).Err()
	if err != nil {
		Log.Errorf("Cannot add user to visit entry. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot add user to visit entry"))
		return
	}

	err = RedisClientVst.Set(id, data, 0).Err()
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func getVisit(w http.ResponseWriter, r *http.Request, id string) {
	Log.Infof("Getting visit with id %s", id)

	result, err := RedisClientVst.Get(id).Result()
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

func updateVisit(w http.ResponseWriter, r *http.Request, id string) {
	Log.Infof("Updaing visit with id %s", id)

	result, err := RedisClientVst.Get(id).Result()
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

	var visit Visit
	err = json.Unmarshal([]byte(result), &visit)
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

	var newVisit Visit
	err = json.Unmarshal(data, &newVisit)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	if newVisit.Location != 0 {
		visit.Location = newVisit.Location
	}
	if newVisit.Mark != 0 {
		visit.Mark = newVisit.Mark
	}
	if newVisit.User != 0 {
		err = RedisClientVstCon.SRem(strconv.Itoa(visit.User), strconv.Itoa(visit.Id)).Err()
		if err != nil {
			Log.Errorf("Cannot remove old user %d from visit with id %d. Reason %s", visit.User, visit.Id, err)
			writeAnswer(w, http.StatusInternalServerError, generateError("Cannot remove old user"))
			return
		}
		err = RedisClientVstCon.SAdd(strconv.Itoa(newVisit.User), strconv.Itoa(visit.Id)).Err()
		if err != nil {
			Log.Errorf("Cannot add new user %d for visit with id %d. Reason %s", newVisit.User, visit.Id, err)
			writeAnswer(w, http.StatusInternalServerError, generateError("Cannot add new user"))
			return
		}
		visit.User = newVisit.User
	}
	if newVisit.VisitedAt != 0 {
		visit.VisitedAt = newVisit.VisitedAt
	}

	resultRedis, err := json.Marshal(visit)
	if err != nil {
		Log.Errorf("Cannot marshal json. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal json"))
		return
	}

	err = RedisClientVst.Set(id, resultRedis, 0).Err()
	if err != nil {
		Log.Errorf("Cannot set redis value. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set redis value"))
		return
	}

	writeAnswer(w, http.StatusOK, "")
}

func processVisit(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	id := path[2]
	if r.Method == "GET" {
		getVisit(w, r, id)
	} else {
		updateVisit(w, r, id)
	}
}
