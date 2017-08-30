package main

import (
	"archive/zip"
	"fmt"
	easyjson "github.com/mailru/easyjson"
	//"github.com/pkg/profile"
	"io/ioutil"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

type JsonFileUsers struct {
	Users []*User `json:"users"`
}

type JsonFileLocations struct {
	Locations []*Location `json:"locations"`
}

type JsonFileVisits struct {
	Visits []*Visit `json:"visits"`
}

/*
func sendRequest(name string, data []byte) error {
	var err error
	switch name {
	case "users":
		user := NewUser(0)
		err = user.UnmarshalJSON(data)
		if err != nil {
			return err
		}
		err = DB.NewUser(user)
		if err != nil {
			return err
		}
	case "locations":
		loc := NewLocation(0)
		err = loc.UnmarshalJSON(data)
		if err != nil {
			return err
		}
		err = DB.NewLocation(loc)
		if err != nil {
			return err
		}
	case "visits":
		visit := NewVisit(0)
		err = visit.UnmarshalJSON(data)
		if err != nil {
			return err
		}
		err = DB.NewVisit(visit)
		if err != nil {
			return err
		}
	default:
		Log.Fatalf("Unknown request type %s", name)
	}
	return nil
}

func processFile(file *zip.File) {
	var structure map[string][]json.RawMessage

	fileReader, err := file.Open()
	if err != nil {
		Log.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	fileData, err := ioutil.ReadAll(fileReader)
	if err != nil {
		Log.Errorf("Cannot read file %s. Reason %s", file.Name, err)
		return
	}

	err = json.Unmarshal(fileData, &structure)
	if err != nil {
		Log.Errorf("Cannot unmarshal file data. Reason %s", err)
		return
	}

	for k, v := range structure {
		for _, raw := range v {
			data, err := raw.MarshalJSON()
			if err != nil {
				Log.Errorf("Cannot marshal raw message. Reason %s", err)
				return
			}
			err = sendRequest(k, data)
			if err != nil {
				Log.Errorf("Cannot send request. Reason %s", err)
				return
			}
		}
	}
}*/

func loadUsers(file *zip.File) {
	fileReader, err := file.Open()
	if err != nil {
		Log.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	data := JsonFileUsers{}
	err = easyjson.UnmarshalFromReader(fileReader, &data)
	if err != nil {
		Log.Errorf("Cannot unmarshal user file. Reason %s", err)
		return
	}

	for _, user := range data.Users {
		err = DB.NewUser(user)
		if err != nil {
			Log.Errorf("Cannot add new user. Reason %s", err)
		}
	}
}

func loadLocations(file *zip.File) {
	fileReader, err := file.Open()
	if err != nil {
		Log.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	data := JsonFileLocations{}
	err = easyjson.UnmarshalFromReader(fileReader, &data)
	if err != nil {
		Log.Errorf("Cannot unmarshal location file. Reason %s", err)
		return
	}

	for _, loc := range data.Locations {
		err = DB.NewLocation(loc)
		if err != nil {
			Log.Errorf("Cannot add new location. Reason %s", err)
		}
	}
}

func loadVisits(file *zip.File) {
	fileReader, err := file.Open()
	if err != nil {
		Log.Errorf("Cannot open file %s. Reason %s", file.Name, err)
		return
	}
	defer fileReader.Close()

	data := JsonFileVisits{}
	err = easyjson.UnmarshalFromReader(fileReader, &data)
	if err != nil {
		Log.Errorf("Cannot unmarshal visit file. Reason %s", err)
		return
	}

	for _, visit := range data.Visits {
		err = DB.NewVisit(visit)
		if err != nil {
			Log.Errorf("Cannot add new visit. Reason %s", err)
		}
	}
}

func load() error {
	reader, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		return fmt.Errorf("Cannot open zip. Reason %s", err)
	}

	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "users") {
			Log.Errorf("Processing file %s", file.Name)
			loadUsers(file)
		} else if strings.HasPrefix(file.Name, "locations") {
			Log.Errorf("Processing file %s", file.Name)
			loadLocations(file)
		} else if strings.HasPrefix(file.Name, "visits") {
			Log.Errorf("Processing file %s", file.Name)
			loadVisits(file)
		} else if file.Name == "options.txt" {
			fileReader, err := file.Open()
			if err != nil {
				Log.Errorf("Cannot open file %s. Reason %s", file.Name, err)
				return nil
			}

			defer fileReader.Close()

			fileData, err := ioutil.ReadAll(fileReader)
			if err != nil {
				Log.Errorf("Cannot read file %s. Reason %s", file.Name, err)
				return nil
			}

			//Log.Errorf("options.txt: %s", string(fileData))
			tsUnix, err := strconv.Atoi(strings.Split(string(fileData), "\n")[0])
			if err != nil {
				ts = time.Now()
				Log.Errorf("Cannot parse ts %s", strings.Split(string(fileData), "\n")[0])
			} else {
				ts = time.Unix(int64(tsUnix), 0)
			}
			//Log.Infof("ts is %s", ts)

		}
	}
	return nil
}

//var profiler prof

func loadToServer() {
	//Log.Infof("Load data to server")

	err := load()
	if err != nil {
		Log.Errorf("Cannot load startup data. Reason %s", err)
	}
	Log.Errorln("ALL LOADED FASTHTTP")

	debug.SetGCPercent(-1)
	//profiler = profile.Start(profile.ProfilePath("."))
}
