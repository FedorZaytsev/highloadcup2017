package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/valyala/fasthttp"
)

var NotFound = fmt.Errorf("Not found")
var CannotParse = fmt.Errorf("Cannot parse")

type Database struct {
	users     map[int]*User
	locations map[int]*Location
	visits    map[int]*Visit
}

func (DB *Database) NewUser(user *User) error {
	//Log.Infof("Inserting user with id %d", user.Id)
	usr, err := DB.GetUser(user.Id)
	if err == nil {
		user.Visits = usr.Visits
	}
	DB.users[user.Id] = user
	return nil
}

func (DB *Database) GetUser(id int) (*User, error) {
	//Log.Infof("Getting user with id %d", id)

	user, ok := DB.users[id]
	if !ok {
		return nil, NotFound
	}

	return user, nil
}

func (DB *Database) UpdateUser(user *User, id int) error {
	//Log.Infof("Updating user with id %d", id)

	return nil
}

func (DB *Database) NewLocation(loc *Location) error {
	//Log.Infof("Inserting location with id %d", loc.Id)
	location, err := DB.GetLocation(loc.Id)
	if err == nil {
		loc.Visits = location.Visits
	}
	DB.locations[loc.Id] = loc
	return nil
}

func (DB *Database) GetLocation(id int) (*Location, error) {
	//Log.Infof("Getting location with id %d", id)

	loc, ok := DB.locations[id]
	if !ok {
		return nil, NotFound
	}

	return loc, nil
}

func (DB *Database) UpdateLocation(loc *Location, id int) error {
	//Log.Infof("Updating location with id %d", id)

	return nil
}

func (DB *Database) NewVisit(visit *Visit) error {
	//Log.Infof("Inserting visit with id %d", visit.Id)

	DB.visits[visit.Id] = visit
	usr, err := DB.GetUser(visit.User)
	if err == NotFound {
		usr = NewUser(visit.User)
		DB.users[usr.Id] = usr
	} else if err != nil {
		return fmt.Errorf("Cannot get user %d. Reason %s", visit.User, err)
	}
	usr.Visits.Add(visit.Id)

	loc, err := DB.GetLocation(visit.Location)
	if err == NotFound {
		loc = NewLocation(visit.Location)
		DB.locations[loc.Id] = loc
	} else if err != nil {
		return fmt.Errorf("Cannot get location %d. Reason %s", visit.Location, err)
	}
	loc.Visits.Add(visit.Id)
	return nil
}

func (DB *Database) GetVisit(id int) (*Visit, error) {
	//Log.Infof("Getting visit with id %d", id)

	v, ok := DB.visits[id]
	if !ok {
		return nil, NotFound
	}
	return v, nil
}

func (DB *Database) UpdateVisit(visit *Visit, id int) error {
	//Log.Infof("Updating visit with id %d", id)

	//DB.visits.Store(visit.Id, visit)
	return nil
}

//select visited_at, mark, place from (select * from visits where id = 1) as v inner join locations on locations.id = v.location where distance < 1000000;
func (DB *Database) GetVisitsFilter(id int, filters *fasthttp.Args) ([]UserVisits, error) {
	result := make([]UserVisits, 0)

	usr, err := DB.GetUser(id)
	if err == NotFound {
		return result, NotFound
	}

	fromDateStr := string(filters.Peek("fromDate"))
	fromDate, err := strconv.Atoi(fromDateStr)
	if err != nil {
		if len(fromDateStr) != 0 {
			return result, CannotParse
		}
		fromDate = math.MinInt32
	}

	toDateStr := string(filters.Peek("toDate"))
	toDate, err := strconv.Atoi(toDateStr)
	if err != nil {
		if toDateStr != "" {
			return result, CannotParse
		}
		toDate = math.MaxInt32
	}

	country := string(filters.Peek("country"))

	toDistance, err := filters.GetUint("toDistance")
	if err == fasthttp.ErrNoArgValue {
		toDistance = math.MaxInt32
	} else if err != nil {
		return result, CannotParse
	}

	usr.Visits.ForEach(func(id int) bool {
		visit, err := DB.GetVisit(id)
		if err == nil && visit.VisitedAt > fromDate && visit.VisitedAt < toDate {
			location, err := DB.GetLocation(visit.Location)
			if err == nil && (country == "" || location.Country == country) && location.Distance < toDistance {
				result = append(result, UserVisits{
					VisitedAt: visit.VisitedAt,
					Mark:      visit.Mark,
					Place:     location.Place,
				})
			}
		}
		return true
	})

	/*DB.visits.Range(func(key, v interface{}) bool {
		visit := v.(*Visit)
		if visit.User == id {
			if visit.VisitedAt > fromDate && visit.VisitedAt < toDate {
				location, err := DB.GetLocation(visit.Location)
				if err == nil && (country == "" || location.Country == country) && location.Distance < toDistance {
					result = append(result, UserVisits{
						VisitedAt: visit.VisitedAt,
						Mark:      visit.Mark,
						Place:     location.Place,
					})
				}
			}
		}
		return true
	})*/

	/*for _, visit := range DB.visits {
		if visit.User == id {
			if visit.VisitedAt > fromDate && visit.VisitedAt < toDate {
				location, err := DB.GetLocation(visit.Location)
				if err == nil && (country == "" || location.Country == country) && location.Distance < toDistance {
					result = append(result, UserVisits{
						VisitedAt: visit.VisitedAt,
						Mark:      visit.Mark,
						Place:     location.Place,
					})
				}
			}
		}
	}*/

	sorter := UserVisitsSorter{
		Data: result,
	}
	sort.Sort(sorter)

	return result, nil
}

//select AVG(mark) from
//(select user, visitedAt, mark from (select * from visits where location=2) as v inner join locations on locations.id = v.location where visitedAt>500) as t inner join users on users.id = t.user where gender = "f";
//
func (DB *Database) GetAverage(id int, filters *fasthttp.Args) (float32, error) {
	var marks float32
	var count int

	loc, err := DB.GetLocation(id)
	if err != nil {
		return 0.0, NotFound
	}

	fromDateStr := string(filters.Peek("fromDate"))
	fromDate, err := strconv.Atoi(fromDateStr)
	if err != nil {
		if fromDateStr != "" {
			return 0.0, CannotParse
		}
		fromDate = math.MinInt32
	}

	toDateStr := string(filters.Peek("toDate"))
	toDate, err := strconv.Atoi(toDateStr)
	if err != nil {
		if toDateStr != "" {
			return 0.0, CannotParse
		}
		toDate = math.MaxInt32
	}

	fromAge, err := filters.GetUint("fromAge")
	if err == fasthttp.ErrNoArgValue {
		fromAge = 0
	} else if err != nil {
		return 0.0, CannotParse
	}

	toAge, err := filters.GetUint("toAge")
	if err == fasthttp.ErrNoArgValue {
		toAge = -1
	} else if err != nil {
		return 0.0, CannotParse
	}

	gender := string(filters.Peek("gender"))
	if gender != "m" && gender != "f" && gender != "" {
		return 0.0, CannotParse
	}

	/*for _, visit := range DB.visits {
		if visit.Location == id {
			if visit.VisitedAt > fromDate && visit.VisitedAt < toDate {
				user, err := DB.GetUser(visit.User)
				if err == nil {
					Log.Warnf("Found user for that visit %#v", user)
					if time.Unix(int64(user.Birthdate), 0).AddDate(fromAge, 0, 0).Before(ts) {
						Log.Warnf("Before ok %v %v", time.Unix(int64(user.Birthdate), 0).AddDate(fromAge, 0, 0), ts)
						if toAge == -1 || time.Unix(int64(user.Birthdate), 0).AddDate(toAge, 0, 0).After(ts) {
							Log.Warnf("Ater ok")
							if gender == "" || user.Gender == gender {
								Log.Infof("Adding %f", float32(visit.Mark))
								marks += float32(visit.Mark)
								count += 1
								Log.Infof("Marks %f %d", marks, count)
							}
						}
					}
				}
			}
		}
	}*/

	loc.Visits.ForEach(func(id int) bool {
		visit, err := DB.GetVisit(id)
		if err == nil && visit.VisitedAt > fromDate && visit.VisitedAt < toDate {
			user, err := DB.GetUser(visit.User)
			if err == nil {
				//Log.Warnf("Found user for that visit %#v", user)
				if time.Unix(int64(user.Birthdate), 0).AddDate(fromAge, 0, 0).Before(ts) {
					//Log.Warnf("Before ok %v %v", time.Unix(int64(user.Birthdate), 0).AddDate(fromAge, 0, 0), ts)
					if toAge == -1 || time.Unix(int64(user.Birthdate), 0).AddDate(toAge, 0, 0).After(ts) {
						//Log.Warnf("Ater ok")
						if gender == "" || user.Gender == gender {
							//Log.Infof("Adding %f", float32(visit.Mark))
							marks += float32(visit.Mark)
							count += 1
							//Log.Infof("Marks %f %d", marks, count)
						}
					}
				}
			}
		}
		return true
	})

	/*DB.visits.Range(func(key, v interface{}) bool {
		visit := v.(*Visit)
		if visit.Location == id {
			if visit.VisitedAt > fromDate && visit.VisitedAt < toDate {
				user, err := DB.GetUser(visit.User)
				if err == nil {
					Log.Warnf("Found user for that visit %#v", user)
					if time.Unix(int64(user.Birthdate), 0).AddDate(fromAge, 0, 0).Before(ts) {
						Log.Warnf("Before ok %v %v", time.Unix(int64(user.Birthdate), 0).AddDate(fromAge, 0, 0), ts)
						if toAge == -1 || time.Unix(int64(user.Birthdate), 0).AddDate(toAge, 0, 0).After(ts) {
							Log.Warnf("Ater ok")
							if gender == "" || user.Gender == gender {
								Log.Infof("Adding %f", float32(visit.Mark))
								marks += float32(visit.Mark)
								count += 1
								Log.Infof("Marks %f %d", marks, count)
							}
						}
					}
				}
			}
		}
		return true
	})*/

	if count == 0 {
		return 0.0, nil
	}

	return marks / float32(count), nil
}

func DatabaseInit() (*Database, error) {
	db := Database{
		users:     make(map[int]*User),
		locations: make(map[int]*Location),
		visits:    make(map[int]*Visit),
	}
	return &db, nil
}
