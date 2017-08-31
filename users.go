package main

import (
	"github.com/valyala/fasthttp"
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

type UserVisitsArray struct {
	Visits []UserVisits `json:"visits"`
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

func NewUser(id int) *User {
	return &User{
		Id:     id,
		Visits: NewArray(),
	}
}

func newUser(ctx *fasthttp.RequestCtx) {
	data := ctx.PostBody()

	body := NewUser(0)
	err := body.UnmarshalJSON(data)
	if err != nil {
		//Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	}
	//Log.Infof("Visits after %p", body.Visits)

	err = DB.NewUser(body)
	if err != nil {
		Log.Errorf("Cannot set id %d. Reason %s", body.Id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot set id"))
		return
	}

	writeAnswer(ctx, http.StatusOK, ANSWER_OK)
}

func getUserVisits(ctx *fasthttp.RequestCtx, id int) {
	filters := ctx.QueryArgs()
	//Log.Infof("Getting user %s visits with filters %#v", id, filters)

	vals, err := DB.GetVisitsFilter(id, filters)
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

	visitsArray := UserVisitsArray{
		Visits: vals,
	}

	result, err := visitsArray.MarshalJSON()
	if err != nil {
		Log.Errorf("Cannot marshal answer. Reason %s", err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot marshal answer"))
		return
	}

	writeAnswer(ctx, http.StatusOK, result)
}

func getUser(ctx *fasthttp.RequestCtx, id int) {
	result, err := DB.GetUser(id)
	if err == NotFound {
		//Log.Infof("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if err != nil {
		Log.Errorf("Cannot get id %d. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("cannot get id"))
		return
	}

	answer, err := result.MarshalJSON()
	if err != nil {
		Log.Errorf("Cannot marshal user %d. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot marshal user"))
		return
	}

	writeAnswer(ctx, http.StatusOK, answer)
}

func updateUser(ctx *fasthttp.RequestCtx, id int) {
	data := ctx.PostBody()

	user, err := DB.GetUser(id)
	if err == NotFound {
		//Log.Infof("Not found")
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if err != nil {
		Log.Errorf("Cannot get user with id %d. Reason %s", id, err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot get user"))
		return
	}

	err = user.UnmarshalJSON(data)
	if err != nil {
		//Log.Warnf("Cannot parse JSON in request. Reason %s", err)
		writeAnswer(ctx, http.StatusBadRequest, []byte{})
		return
	}

	err = DB.UpdateUser(user, id)
	if err != nil {
		Log.Errorf("Cannot update user. Reason %s", err)
		writeAnswer(ctx, http.StatusInternalServerError, generateError("Cannot update user"))
		return
	}

	writeAnswer(ctx, http.StatusOK, ANSWER_OK)
}

func processUser(ctx *fasthttp.RequestCtx) {
	path := strings.Split(string(ctx.Path()), "/")
	id, err := strconv.Atoi(path[2])
	if err != nil {
		//Log.Infof("Cannot parse id %s. Reason %s", path[2], err)
		writeAnswer(ctx, http.StatusNotFound, []byte{})
		return
	}
	if len(path) == 4 {
		getUserVisits(ctx, id)
		return
	}
	if string(ctx.Method()) == "GET" {
		getUser(ctx, id)
	} else {
		updateUser(ctx, id)
	}
}
