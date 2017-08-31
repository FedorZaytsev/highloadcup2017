package main

import (
	"github.com/valyala/fasthttp"
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

type LocationAvg struct {
	Avg float64 `json:"avg"`
}

func NewLocation(id int) *Location {
	return &Location{
		Id:     id,
		Visits: NewArray(),
	}
}

func newLocation(ctx *fasthttp.RequestCtx) {
	data := ctx.PostBody()

	body := NewLocation(0)
	err := body.UnmarshalJSON(data)
	if err != nil {
		//Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	}

	err = DB.NewLocation(body)
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(ctx, http.StatusOK, ANSWER_OK)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func getLocationAvg(ctx *fasthttp.RequestCtx, id int) {
	filters := ctx.QueryArgs()
	//Log.Infof("Getting avg for location %d with filters %#v", id, filters)

	val, err := DB.GetAverage(id, filters)
	switch err {
	case NotFound:
		//Log.Infof("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	case CannotParse:
		//Log.Infof("Cannot parse")
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	case nil:
		break
	default:
		Log.Errorf("Error while getting user visits. Reason %s", err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Error while getting user visits"))
		return
	}

	avg := LocationAvg{
		Avg: toFixed(float64(val), 5),
	}

	result, err := avg.MarshalJSON()
	if err != nil {
		Log.Errorf("Cannot marshal answer. Reason %s", err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot marshal answer"))
		return
	}

	writeAnswer(ctx, http.StatusOK, result)
}

func getLocation(ctx *fasthttp.RequestCtx, id int) {
	location, err := DB.GetLocation(id)
	if err == NotFound {
		//Log.Warnf("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if err != nil {
		//Log.Errorf("Cannot get id %s. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	result, err := location.MarshalJSON()
	if err != nil {
		Log.Errorf("Cannot marshal location %#v. Reason %s", location, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot marshal location"))
		return
	}

	writeAnswer(ctx, http.StatusOK, result)
}

func updateLocation(ctx *fasthttp.RequestCtx, id int) {
	data := ctx.PostBody()

	loc, err := DB.GetLocation(id)
	if err == NotFound {
		//Log.Infof("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if err != nil {
		Log.Errorf("Cannot get location with id %d. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot get location"))
		return
	}

	err = loc.UnmarshalJSON(data)
	if err != nil {
		//Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	}

	err = DB.UpdateLocation(loc, id)
	if err != nil {
		Log.Errorf("Cannot update location. Reason %s", err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot update location"))
		return
	}

	writeAnswer(ctx, http.StatusOK, ANSWER_OK)
}

func processLocation(ctx *fasthttp.RequestCtx) {
	path := strings.Split(string(ctx.Path()), "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		//Log.Infof("Cannot parse id %s. Reason %s", path[2], err)
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if len(path) == 4 {
		getLocationAvg(ctx, id)
		return
	}
	if string(ctx.Method()) == "GET" {
		getLocation(ctx, id)
	} else {
		updateLocation(ctx, id)
	}
}
