package main

import (
	"errors"
	"io/ioutil"
	"log/syslog"
	"os"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

//Log is logger for all programm
var Log *log.Logger

//initLogger take from Config parameters for logger and init logger.
//If we use syslog, we will call initSyslogger.
func initLogger() (*log.Logger, error) {
	logger := log.New()
	switch Config.Log.Type {
	case "syslog":
		logger = initSyslogger()
	case "stderr":
		logger.Out = os.Stderr
	case "stdout":
		logger.Out = os.Stdout
	default:
		return nil, errors.New("Incorrect out type for log.")
	}
	if Config.Log.DebugMode {
		logger.Out = os.Stdout
	}
	logger.Formatter = &log.TextFormatter{}
	logger.Level = logLevel[Config.Log.Severity]
	return logger, nil
}

var logLevel = map[string]log.Level{
	"LOG_EMERG":   log.PanicLevel,
	"LOG_ALERT":   log.PanicLevel,
	"LOG_CRIT":    log.FatalLevel,
	"LOG_ERR":     log.ErrorLevel,
	"LOG_WARNING": log.WarnLevel,
	"LOG_NOTICE":  log.InfoLevel,
	"LOG_INFO":    log.InfoLevel,
	"LOG_DEBUG":   log.DebugLevel,
}

//initSyslogger connectes to SysLog and init Logger for it
func initSyslogger() *log.Logger {
	var LogSeverity = map[string]syslog.Priority{
		"LOG_EMERG":   syslog.LOG_EMERG,
		"LOG_ALERT":   syslog.LOG_ALERT,
		"LOG_CRIT":    syslog.LOG_CRIT,
		"LOG_ERR":     syslog.LOG_ERR,
		"LOG_WARNING": syslog.LOG_WARNING,
		"LOG_NOTICE":  syslog.LOG_NOTICE,
		"LOG_INFO":    syslog.LOG_INFO,
		"LOG_DEBUG":   syslog.LOG_DEBUG,
	}
	var LogFacility = map[string]syslog.Priority{
		"LOG_KERN":     syslog.LOG_KERN,
		"LOG_USER":     syslog.LOG_USER,
		"LOG_MAIL":     syslog.LOG_MAIL,
		"LOG_DAEMON":   syslog.LOG_DAEMON,
		"LOG_AUTH":     syslog.LOG_AUTH,
		"LOG_SYSLOG":   syslog.LOG_SYSLOG,
		"LOG_LPR":      syslog.LOG_LPR,
		"LOG_NEWS":     syslog.LOG_NEWS,
		"LOG_UUCP":     syslog.LOG_UUCP,
		"LOG_CRON":     syslog.LOG_CRON,
		"LOG_AUTHPRIV": syslog.LOG_AUTHPRIV,
		"LOG_FTP":      syslog.LOG_FTP,

		"LOG_LOCAL0": syslog.LOG_LOCAL0,
		"LOG_LOCAL1": syslog.LOG_LOCAL1,
		"LOG_LOCAL2": syslog.LOG_LOCAL2,
		"LOG_LOCAL3": syslog.LOG_LOCAL3,
		"LOG_LOCAL4": syslog.LOG_LOCAL4,
		"LOG_LOCAL5": syslog.LOG_LOCAL5,
		"LOG_LOCAL6": syslog.LOG_LOCAL6,
		"LOG_LOCAL7": syslog.LOG_LOCAL7,
	}
	logger := log.New()
	hook, err := logrus_syslog.NewSyslogHook(
		Config.Log.NetworkType,
		Config.Log.Host+":"+Config.Log.Port,
		LogSeverity[Config.Log.Severity]|LogFacility[Config.Log.Facility],
		Config.Title)
	if err != nil {
		log.Errorln(err)
		hook, _ = logrus_syslog.NewSyslogHook(
			"", "", LogSeverity[Config.Log.Severity]|LogFacility[Config.Log.Facility],
			Config.Title)
	}
	logger.Hooks.Add(hook)

	logger.Out = ioutil.Discard
	return logger
}
