package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/ogier/pflag"
)

var (
	Version           string
	BuildTime         string
	RedisClientUsr    *redis.Client
	RedisClientLoc    *redis.Client
	RedisClientVst    *redis.Client
	RedisClientVstCon *redis.Client
)

func init() {
	var versReq bool
	pflag.StringVarP(&configPath, "config", "c", "config.toml", "Used for set path to config file.")
	pflag.BoolVarP(&versReq, "version", "v", false, "Use for build time and version print")
	var err error
	pflag.Parse()
	if versReq {
		fmt.Println("Version: ", Version)
		fmt.Println("Build time:", BuildTime)
		os.Exit(0)
	}
	Config, err = configure()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	Log, err = initLogger()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	RedisClientUsr = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	err = RedisClientUsr.Ping().Err()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	RedisClientLoc = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	err = RedisClientLoc.Ping().Err()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	RedisClientVst = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	})
	err = RedisClientVst.Ping().Err()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}

	RedisClientVstCon = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       3,
	})
	err = RedisClientVstCon.Ping().Err()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}

func writeAnswer(w http.ResponseWriter, code int, body string) {
	w.Header().Set("Cache-Control", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, body)
}

func generateError(data string) string {
	return fmt.Sprintf("{\"error\": \"%s\"}", data)
}

func main() {
	Log.Infof("Started\n")

	http.HandleFunc("/users/new", newUser)
	http.HandleFunc("/users/", processUser)
	http.HandleFunc("/locations/new", newLocation)
	http.HandleFunc("/locations/", processLocation)
	http.HandleFunc("/visits/new", newVisit)
	http.HandleFunc("/visits/", processVisit)
	http.ListenAndServe(":8080", nil)

	//uncomment if it is a demon
	//sgnl := make(chan os.Signal, 1)
	//signal.Notify(sgnl, os.Interrupt, syscall.SIGTERM)
	//<-sgnl
}
