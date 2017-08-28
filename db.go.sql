package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

var NotFound = fmt.Errorf("Not found")

type Database struct {
	sql *sql.DB
}

func (DB *Database) Init() error {
	Log.Infof("Initing db")
	_, err := DB.sql.Exec("DROP TABLE IF EXISTS Users")
	if err != nil {
		return err
	}
	_, err = DB.sql.Exec("DROP TABLE IF EXISTS Locations")
	if err != nil {
		return err
	}
	_, err = DB.sql.Exec("DROP TABLE IF EXISTS Visits")
	if err != nil {
		return err
	}
	_, err = DB.sql.Exec("CREATE TABLE IF NOT EXISTS Users (id INT NOT NULL PRIMARY KEY, email CHAR(100), firstname CHAR(50), lastname CHAR(50), gender CHAR(1), birthdate INT)")
	if err != nil {
		return err
	}
	_, err = DB.sql.Exec("CREATE TABLE IF NOT EXISTS Locations (id INT NOT NULL PRIMARY KEY, place TEXT, country CHAR(50), city CHAR(50), distance INT)")
	if err != nil {
		return err
	}
	_, err = DB.sql.Exec("CREATE TABLE IF NOT EXISTS Visits (id INT NOT NULL PRIMARY KEY, location INT NOT NULL, user INT NOT NULL, visitedAt INT, mark INT)")
	if err != nil {
		return err
	}
	Log.Infof("DB init done")
	return nil
}

func (DB *Database) NewUser(user User) error {
	Log.Infof("Inserting user with id %d", user.Id)
	_, err := DB.sql.Exec("INSERT INTO Users VALUES (?, ?, ?, ?, ?, ?)", user.Id, user.Email, user.FirstName, user.LastName, user.Gender, user.Birthdate)
	if err != nil {
		return fmt.Errorf("Cannot insert new user. Reason %s", err)
	}
	return nil
}

func (DB *Database) GetUser(id int) (User, error) {
	Log.Infof("Getting user with id %d", id)

	var user User
	val := DB.sql.QueryRow("SELECT * from Users WHERE id = ?", id)
	err := val.Scan(&user.Id, &user.Email, &user.FirstName, &user.LastName, &user.Gender, &user.Birthdate)
	if err == sql.ErrNoRows {
		return user, NotFound
	}
	return user, err
}

func (DB *Database) UpdateUser(user User, id int) error {
	Log.Infof("Updating user with id %d", id)

	rowsTemp, err := DB.sql.Query("select * from Users where id=?", id)
	if err != nil {
		return fmt.Errorf("Cannot check is user exists. Reason %s", err)
	}
	defer rowsTemp.Close()
	if !rowsTemp.Next() {
		return NotFound
	}

	v := reflect.ValueOf(user)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			//Log.Infof("upd %s", fmt.Sprintf("UPDATE Users SET %s=%v where id=%d", v.Type().Field(i).Name, field.Interface(), id))
			_, err := DB.sql.Exec(fmt.Sprintf("UPDATE Users SET %s=? where id=?", v.Type().Field(i).Name), field.Interface(), id)
			if err != nil {
				return fmt.Errorf("Cannot update field %s for object %#v. Reason %s", v.Type().Field(i).Name, user, err)
			}
		}
	}
	return nil
}

func (DB *Database) NewLocation(loc Location) error {
	Log.Infof("Inserting location with id %d", loc.Id)
	_, err := DB.sql.Exec("INSERT INTO Locations VALUES (?, ?, ?, ?, ?)", loc.Id, loc.Place, loc.Country, loc.City, loc.Distance)
	if err != nil {
		return fmt.Errorf("Cannot insert new location. Reason %s", err)
	}
	return nil
}

func (DB *Database) GetLocation(id int) (Location, error) {
	Log.Infof("Getting location with id %d", id)

	var loc Location
	val := DB.sql.QueryRow("SELECT * from Locations WHERE id = ?", id)
	err := val.Scan(&loc.Id, &loc.Place, &loc.Country, &loc.City, &loc.Distance)
	if err == sql.ErrNoRows {
		return loc, NotFound
	}
	return loc, err
}

func (DB *Database) UpdateLocation(loc Location, id int) error {
	Log.Infof("Updating location with id %d", id)

	rowsTemp, err := DB.sql.Query("select * from Locations where id=?", id)
	if err != nil {
		return fmt.Errorf("Cannot check is location exists. Reason %s", err)
	}
	defer rowsTemp.Close()
	if !rowsTemp.Next() {
		return NotFound
	}

	v := reflect.ValueOf(loc)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			//Log.Infof("upd %s", fmt.Sprintf("UPDATE Locations SET %s=%v where id=%d", v.Type().Field(i).Name, field.Interface(), id))
			_, err := DB.sql.Exec(fmt.Sprintf("UPDATE Locations SET %s=? where id=?", v.Type().Field(i).Name), field.Interface(), id)
			if err != nil {
				return fmt.Errorf("Cannot update field %s for object %#v. Reason %s", v.Type().Field(i).Name, loc, err)
			}
		}
	}
	return nil
}

func (DB *Database) NewVisit(visit Visit) error {
	Log.Infof("Inserting visit with id %d", visit.Id)
	_, err := DB.sql.Exec("INSERT INTO Visits VALUES (?, ?, ?, ?, ?)", visit.Id, visit.Location, visit.User, visit.VisitedAt, visit.Mark)
	if err != nil {
		return fmt.Errorf("Cannot insert new visit. Reason %s", err)
	}
	return nil
}

func (DB *Database) GetVisit(id int) (Visit, error) {
	Log.Infof("Getting visit with id %d", id)

	var v Visit
	val := DB.sql.QueryRow("SELECT * from Visits WHERE id = ?", id)
	err := val.Scan(&v.Id, &v.Location, &v.User, &v.VisitedAt, &v.Mark)
	if err == sql.ErrNoRows {
		return v, NotFound
	}
	return v, err
}

func (DB *Database) UpdateVisit(visit Visit, id int) error {
	Log.Infof("Updating visit with id %d", id)

	rowsTemp, err := DB.sql.Query("select * from Visits where id=?", id)
	if err != nil {
		return fmt.Errorf("Cannot check is visit exists. Reason %s", err)
	}
	defer rowsTemp.Close()
	if !rowsTemp.Next() {
		return NotFound
	}

	v := reflect.ValueOf(visit)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Interface() != reflect.Zero(field.Type()).Interface() {
			//Log.Infof("upd %s", fmt.Sprintf("UPDATE Visits SET %s=%v where id=%d", v.Type().Field(i).Name, field.Interface(), id))
			_, err := DB.sql.Exec(fmt.Sprintf("UPDATE Visits SET %s=? where id=?", v.Type().Field(i).Name), field.Interface(), id)
			if err != nil {
				return fmt.Errorf("Cannot update field %s for object %#v. Reason %s", v.Type().Field(i).Name, visit, err)
			}
		}
	}
	return nil
}

//select visited_at, mark, place from (select * from visits where id = 1) as v inner join locations on locations.id = v.location where distance < 1000000;
func (DB *Database) GetVisitsFilter(id int, filters url.Values) ([]UserVisits, error) {
	result := make([]UserVisits, 0)

	rowsTemp, err := DB.sql.Query("select * from Users where id=?", id)
	if err != nil {
		return result, fmt.Errorf("Cannot check is user exists. Reason %s", err)
	}
	defer rowsTemp.Close()
	if !rowsTemp.Next() {
		return result, NotFound
	}

	where, err := generateWhereClasure(filters)
	if err != nil {
		return result, err
	}
	Log.Infof("select visitedAt, mark, place from (select * from Visits where user = " + strconv.Itoa(id) + ") as v inner join Locations on locations.id = v.location " + where)
	rows, err := DB.sql.Query(`
select visitedAt, mark, place from (
	select * from Visits where user = ?
) as v inner join Locations on Locations.id = v.location `+where, id)

	if err != nil {
		return result, fmt.Errorf("Cannot execute select for vistis filter. Reason %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var v UserVisits
		err := rows.Scan(&v.VisitedAt, &v.Mark, &v.Place)
		Log.Infof("row next %v", v)
		if err != nil {
			return result, fmt.Errorf("Cannot get user visit. Reason %s", err)
		}
		result = append(result, v)
	}
	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("Error while executing rows.Next(). Reason %s", err)
	}
	return result, nil
}

//select AVG(mark) from
//(select user, visitedAt, mark from (select * from visits where location=2) as v inner join locations on locations.id = v.location where visitedAt>500) as t inner join users on users.id = t.user where gender = "f";
//
func (DB *Database) GetAverage(id int, filters url.Values) (float32, error) {

	rowsTemp, err := DB.sql.Query("select * from Locations where id=?", id)
	if err != nil {
		return 0.0, fmt.Errorf("Cannot check is location exists. Reason %s", err)
	}
	defer rowsTemp.Close()
	if !rowsTemp.Next() {
		return 0.0, NotFound
	}

	inner, err := generateWhereClasureAvgInner(filters)
	if err != nil {
		return 0.0, err
	}
	outter, err := generateWhereClasureAvgOutter(filters)
	if err != nil {
		return 0.0, err
	}

	req := `
select SUM(mark), COUNT(mark) from (
	select user, visitedAt, mark from (
		select * from Visits where location=?
	) as v inner join Locations on Locations.id = v.location ` + inner + `
) as t inner join Users on Users.id = t.user ` + outter

	var sum, count float32
	val := DB.sql.QueryRow(req, id)
	err = val.Scan(&sum, &count)
	if err == sql.ErrNoRows {
		return 0.0, NotFound
	}
	if math.Abs(float64(count)) < 1e-5 {
		return 0.0, nil
	}
	Log.Infof("sum %f count %f", sum, count)

	return sum / count, nil
}

func generateWhereClasure(filters url.Values) (string, error) {
	var buf bytes.Buffer
	if filters.Get("fromDate") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		if _, err := strconv.Atoi(filters.Get("fromDate")); err != nil {
			return "", fmt.Errorf("Cannot convert fromDate %s", filters.Get("fromDate"))
		}
		buf.WriteString("visitedAt > " + filters.Get("fromDate"))
	}
	if filters.Get("toDate") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		if _, err := strconv.Atoi(filters.Get("toDate")); err != nil {
			return "", fmt.Errorf("Cannot convert toDate %s", filters.Get("toDate"))
		}
		buf.WriteString("visitedAt < " + filters.Get("toDate"))
	}
	if filters.Get("country") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		buf.WriteString("country = \"" + filters.Get("country") + "\"")
	}
	if filters.Get("toDistance") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		if _, err := strconv.Atoi(filters.Get("toDistance")); err != nil {
			return "", fmt.Errorf("Cannot convert toDistance %s", filters.Get("toDistance"))
		}
		buf.WriteString("distance < " + filters.Get("toDistance"))
	}

	if buf.Len() > 0 {
		return "WHERE " + buf.String(), nil
	}
	return "", nil
}

func generateWhereClasureAvgInner(filters url.Values) (string, error) {
	var buf bytes.Buffer
	if filters.Get("fromDate") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		if _, err := strconv.Atoi(filters.Get("fromDate")); err != nil {
			return "", fmt.Errorf("Cannot convert fromDate %s", filters.Get("fromDate"))
		}
		buf.WriteString("visitedAt > " + filters.Get("fromDate"))
	}
	if filters.Get("toDate") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		if _, err := strconv.Atoi(filters.Get("toDate")); err != nil {
			return "", fmt.Errorf("Cannot convert toDate %s", filters.Get("toDate"))
		}
		buf.WriteString("visitedAt < " + filters.Get("toDate"))
	}

	if buf.Len() > 0 {
		return "WHERE " + buf.String(), nil
	}
	return "", nil
}

func generateWhereClasureAvgOutter(filters url.Values) (string, error) {
	var buf bytes.Buffer
	if filters.Get("fromAge") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		fromAge, err := strconv.Atoi(filters.Get("fromAge"))
		if err != nil {
			return "", fmt.Errorf("Cannot parse fromAge. %s Reason %s", filters.Get("fromAge"), err)
		}
		fromDateAge := time.Unix(0, 0).AddDate(fromAge, 0, 0).Unix()
		buf.WriteString("birthdate + " + strconv.FormatInt(fromDateAge, 10) + " < " + strconv.FormatInt(time.Now().Unix(), 10))
	}
	if filters.Get("toAge") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		toAge, err := strconv.Atoi(filters.Get("toAge"))
		if err != nil {
			return "", fmt.Errorf("Cannot parse toAge. %s Reason %s", filters.Get("toAge"), err)
		}
		toDateAge := time.Unix(0, 0).AddDate(toAge, 0, 0).Unix()
		buf.WriteString("birthdate + " + strconv.FormatInt(toDateAge, 10) + " > " + strconv.FormatInt(time.Now().Unix(), 10))
	}
	if filters.Get("gender") != "" {
		if buf.Len() != 0 {
			buf.WriteString(" and ")
		}
		buf.WriteString("gender = \"" + filters.Get("gender") + "\"")
	}

	if buf.Len() > 0 {
		return "WHERE " + buf.String(), nil
	}
	return "", nil
}

func DatabaseInit() (*Database, error) {
	mysql, err := sql.Open("mysql", "root:@/highloadcup?parseTime=true&charset=utf8&collation=utf8_general_ci")
	if err != nil {
		return nil, fmt.Errorf("Cannot create connection to mysql. Reason %s", err)
	}
	db := Database{
		sql: mysql,
	}
	err = db.Init()
	if err != nil {
		return nil, fmt.Errorf("Cannot init DB. Reason %s", err)
	}
	return &db, nil
}
