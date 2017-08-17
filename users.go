package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type User struct {
	Id        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Gender    string `json:"gender"`
	Birthdate int64  `json:"birth_date"`
}

func newUser(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Log.Warnf("Cannot read body from request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot read body from request"))
		return
	}

	var body User
	err = json.Unmarshal(data, &body)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	err = RedisClientUsr.Set(strconv.Itoa(body.Id), data, 0).Err()
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(w, http.StatusOK, "{}")
}

func getUserVisits(w http.ResponseWriter, r *http.Request, id string) {
}

func getUser(w http.ResponseWriter, r *http.Request, id string) {
	Log.Infof("Getting user with id %s", id)

	result, err := RedisClientUsr.Get(id).Result()
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

func updateUser(w http.ResponseWriter, r *http.Request, id string) {
	Log.Infof("Updaing user with id %s", id)

	result, err := RedisClientUsr.Get(id).Result()
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

	var user User
	err = json.Unmarshal([]byte(result), &user)
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

	var newUser User
	err = json.Unmarshal(data, &newUser)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(w, http.StatusBadRequest, generateError("Cannot parse JSON in request"))
		return
	}

	if newUser.Birthdate != 0 {
		user.Birthdate = newUser.Birthdate
	}
	if newUser.Email != "" {
		user.Email = newUser.Email
	}
	if newUser.FirstName != "" {
		user.FirstName = newUser.FirstName
	}
	if newUser.Gender != "" {
		user.Gender = newUser.Gender
	}
	if newUser.LastName != "" {
		user.LastName = newUser.LastName
	}

	resultRedis, err := json.Marshal(user)
	if err != nil {
		Log.Errorf("Cannot marshal json. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot marshal json"))
		return
	}

	err = RedisClientUsr.Set(id, resultRedis, 0).Err()
	if err != nil {
		Log.Errorf("Cannot set redis value. Reason %s", err)
		writeAnswer(w, http.StatusInternalServerError, generateError("Cannot set redis value"))
		return
	}

	writeAnswer(w, http.StatusOK, "")
}

func processUser(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")
	if path[len(path)-1] == "visits" {
		getUserVisits(w, r, path[len(path)-2])
		return
	}
	if r.Method == "GET" {
		getUser(w, r, path[len(path)-1])
	} else {
		updateUser(w, r, path[len(path)-1])
	}
}
