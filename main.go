package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ogier/pflag"
	//"github.com/pkg/profile"
	"github.com/valyala/fasthttp"
)

var (
	Version   string
	BuildTime string
	ts        time.Time
	DB        *Database
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

	DB, err = DatabaseInit()
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	ts = time.Now()
}

func writeAnswer(ctx *fasthttp.RequestCtx, code int, body string) {
	ctx.SetContentType("application/json; charset=UTF-8")
	ctx.SetStatusCode(code)
	fmt.Fprintf(ctx, body)
}

func generateError(data string) string {
	return fmt.Sprintf("{\"error\": \"%s\"}", data)
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	Log.Infof("PATH %s", string(path))
	switch string(path[:2]) {
	case "/u":
		if path[len(path)-1] == 'w' {
			newUser(ctx)
		} else {
			processUser(ctx)
		}
	case "/l":
		if path[len(path)-1] == 'w' {
			newLocation(ctx)
		} else {
			processLocation(ctx)
		}
	case "/v":
		if path[len(path)-1] == 'w' {
			newVisit(ctx)
		} else {
			processVisit(ctx)
		}
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func main() {
	Log.Infof("Started\n")
	//defer profile.Start(profile.ProfilePath(".")).Stop()

	go loadToServer()

	err := fasthttp.ListenAndServe(":8080", requestHandler)
	if err != nil {
		Log.Errorf("ListAndServe error: %s", err)
	}

	/*http.HandleFunc("/users/new", newUser)
	http.HandleFunc("/users/", processUser)
	http.HandleFunc("/locations/new", newLocation)
	http.HandleFunc("/locations/", processLocation)
	http.HandleFunc("/visits/new", newVisit)
	http.HandleFunc("/visits/", processVisit)
	http.ListenAndServe(":80", nil)*/

	//uncomment if it is a demon
	//sgnl := make(chan os.Signal, 1)
	//signal.Notify(sgnl, os.Interrupt, syscall.SIGTERM)
	//<-sgnl
}
