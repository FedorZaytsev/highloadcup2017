package main

import (
	"github.com/valyala/fasthttp"
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

func NewVisit(id int) *Visit {
	return &Visit{
		Id: id,
	}
}

var ANSWER_OK = []byte("{}")

func newVisit(ctx *fasthttp.RequestCtx) {
	data := ctx.PostBody()

	var body Visit
	err := body.UnmarshalJSON(data)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	}

	err = DB.NewVisit(&body)
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(ctx, http.StatusOK, ANSWER_OK)
}

func getVisit(ctx *fasthttp.RequestCtx, id int) {
	visit, err := DB.GetVisit(id)
	if err == NotFound {
		Log.Warnf("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if err != nil {
		Log.Errorf("Cannot get id %s. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	result, err := visit.MarshalJSON()
	if err != nil {
		Log.Errorf("Cannot marshal visit %#v. Reason %s", visit, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot marshal visit"))
		return
	}

	writeAnswer(ctx, http.StatusOK, result)
}

func updateVisit(ctx *fasthttp.RequestCtx, id int) {
	data := ctx.PostBody()

	visit, err := DB.GetVisit(id)
	if err == NotFound {
		Log.Infof("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if err != nil {
		Log.Errorf("Cannot get visit with id %d. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot get visit"))
		return
	}

	oldUser := visit.User
	oldLocation := visit.Location

	err = visit.UnmarshalJSON(data)
	if err != nil {
		Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	}

	//UPDATE USER
	if visit.User != oldUser {
		usrOld, err := DB.GetUser(oldUser)
		if err != nil {
			Log.Errorf("Cannot get user %d. Reason %s", oldUser, err)
			writeAnswer(ctx, http.StatusBadRequest, []byte{})
			return
		}
		usrOld.Visits.Remove(visit.Id)

		usr, err := DB.GetUser(visit.User)
		if err != nil {
			Log.Errorf("Cannot get user %d. Reason %s", visit.User, err)
			writeAnswer(ctx, http.StatusBadRequest, []byte{})
			return
		}
		usr.Visits.Add(visit.Id)
	}

	//UPDATE LOCATION
	if visit.Location != oldLocation {
		locOld, err := DB.GetLocation(oldLocation)
		if err != nil {
			Log.Errorf("Cannot get location %d. Reason %s", oldLocation, err)
			writeAnswer(ctx, http.StatusBadRequest, []byte{})
			return
		}
		locOld.Visits.Remove(visit.Id)

		loc, err := DB.GetLocation(visit.Location)
		if err != nil {
			Log.Errorf("Cannot get location %d. Reason %s", visit.Location, err)
			writeAnswer(ctx, http.StatusBadRequest, []byte{})
			return
		}
		loc.Visits.Add(visit.Id)
	}

	err = DB.UpdateVisit(visit, id)
	if err != nil {
		Log.Errorf("Cannot update visit. Reason %s", err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot update visit"))
		return
	}

	writeAnswer(ctx, http.StatusOK, ANSWER_OK)
}

func processVisit(ctx *fasthttp.RequestCtx) {
	path := strings.Split(string(ctx.Path()), "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		Log.Infof("Cannot parse id %s. Reason %s", path[2], err)
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if string(ctx.Method()) == "GET" {
		getVisit(ctx, id)
	} else {
		updateVisit(ctx, id)
	}
}
