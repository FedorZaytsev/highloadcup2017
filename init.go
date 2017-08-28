package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func sendRequest(name, data string) error {
	var err error
	switch name {
	case "users":
		user := User{
			Visits: NewArray(),
		}
		err = json.Unmarshal([]byte(data), &user)
		if err != nil {
			return err
		}
		err = DB.NewUser(&user)
		if err != nil {
			return err
		}
	case "locations":
		loc := Location{
			Visits: NewArray(),
		}
		err = json.Unmarshal([]byte(data), &loc)
		if err != nil {
			return err
		}
		err = DB.NewLocation(&loc)
		if err != nil {
			return err
		}
	case "visits":
		var visit Visit
		err = json.Unmarshal([]byte(data), &visit)
		if err != nil {
			return err
		}
		err = DB.NewVisit(&visit)
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
			err = sendRequest(k, string(data))
			if err != nil {
				Log.Errorf("Cannot send request. Reason %s", err)
				return
			}
		}
	}
}

func load() error {
	reader, err := zip.OpenReader("/tmp/data/data.zip")
	if err != nil {
		return fmt.Errorf("Cannot open zip. Reason %s", err)
	}

	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "locations") || strings.HasPrefix(file.Name, "users") || strings.HasPrefix(file.Name, "visits") {
			Log.Infof("Processing file %s", file.Name)

			processFile(file)
		}
		if file.Name == "options.txt" {
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

			Log.Errorf("options.txt: %s", string(fileData))
			tsUnix, err := strconv.Atoi(strings.Split(string(fileData), "\n")[0])
			if err != nil {
				ts = time.Now()
				Log.Errorf("Cannot parse ts %s", strings.Split(string(fileData), "\n")[0])
			} else {
				ts = time.Unix(int64(tsUnix), 0)
			}
			Log.Infof("ts is %s", ts)

		}
	}
	return nil
}

func loadToServer() {
	Log.Infof("Load data to server")

	err := load()
	if err != nil {
		Log.Errorf("Cannot load startup data. Reason %s", err)
	}
	Log.Errorln("ALL LOADED")
}
