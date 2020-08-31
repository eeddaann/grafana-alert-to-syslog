package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	//"log/syslog"
	"net/http"
	"strings"

	syslog "github.com/RackSec/srslog"
)

type Config struct {
	SyslogHost string
	Protocol   string
	Port       string
}

var config Config

// RawAlert is alert from Grafana
type RawAlert struct {
	RuleName string
	Tags     Tags
	Title    string
	Message  string
	State    string
}

// Tags are the tags of the alert from Grafana
type Tags struct {
	Tag      string
	Priority string
}

// FormattedLog is log in the format of syslog
type FormattedLog struct {
	Msg      string
	Priority syslog.Priority
	Tag      string
}

func formatLog(alert RawAlert) FormattedLog {
	formatted := FormattedLog{Priority: syslog.LOG_ALERT, Tag: "Grafana"} // set default values
	formatted.Msg = alert.Title                                           // set log message as the title since it's include the state - not sure about this
	if alert.Tags.Priority != "" {
		switch strings.ToUpper(alert.Tags.Priority) {
		// set severity
		case "EMERG":
			formatted.Priority = syslog.LOG_EMERG
		case "ALERT":
			formatted.Priority = syslog.LOG_ALERT
		case "CRIT":
			formatted.Priority = syslog.LOG_CRIT
		case "ERR":
			formatted.Priority = syslog.LOG_ERR
		case "WARNING":
			formatted.Priority = syslog.LOG_WARNING
		case "NOTICE":
			formatted.Priority = syslog.LOG_NOTICE
		case "INFO":
			formatted.Priority = syslog.LOG_INFO
		case "DEBUG":
			formatted.Priority = syslog.LOG_DEBUG
		case "KERN":
			formatted.Priority = syslog.LOG_KERN
		case "USER":
			formatted.Priority = syslog.LOG_USER
		case "MAIL":
			formatted.Priority = syslog.LOG_MAIL
		case "DAEMON":
			formatted.Priority = syslog.LOG_DAEMON
		case "AUTH":
			formatted.Priority = syslog.LOG_AUTH
		case "SYSLOG":
			formatted.Priority = syslog.LOG_SYSLOG
		case "LPR":
			formatted.Priority = syslog.LOG_LPR
		case "NEWS":
			formatted.Priority = syslog.LOG_NEWS
		case "UUCP":
			formatted.Priority = syslog.LOG_UUCP
		case "CRON":
			formatted.Priority = syslog.LOG_CRON
		case "AUTHPRIV":
			formatted.Priority = syslog.LOG_AUTHPRIV
		case "FTP":
			formatted.Priority = syslog.LOG_FTP
		case "LOCAL0":
			formatted.Priority = syslog.LOG_LOCAL0
		case "LOCAL1":
			formatted.Priority = syslog.LOG_LOCAL1
		case "LOCAL2":
			formatted.Priority = syslog.LOG_LOCAL2
		case "LOCAL3":
			formatted.Priority = syslog.LOG_LOCAL3
		case "LOCAL4":
			formatted.Priority = syslog.LOG_LOCAL4
		case "LOCAL5":
			formatted.Priority = syslog.LOG_LOCAL5
		case "LOCAL6":
			formatted.Priority = syslog.LOG_LOCAL6
		case "LOCAL7":
			formatted.Priority = syslog.LOG_LOCAL7
		default:
			log.Printf("could not set priority %v (case insensitive)", alert.Tags.Priority)
			formatted.Priority = syslog.LOG_ALERT
		}
	}
	if alert.Tags.Tag != "" {
		// set custom tag
		formatted.Tag = alert.Tags.Tag
	}
	return formatted
}
func parseAlert(alert string) RawAlert {
	var alertObj RawAlert
	if err := json.Unmarshal([]byte(alert), &alertObj); err != nil {
		panic(err)
	}
	return alertObj
}

func status(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "up and running!")
}

func handle(w http.ResponseWriter, req *http.Request) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
	}
	println(string(b))
	var alertObj = parseAlert(string(b))
	var formattedLog = formatLog(alertObj)
	sysLog, err := syslog.Dial(config.Protocol, config.SyslogHost,
		formattedLog.Priority|syslog.LOG_USER, formattedLog.Tag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(sysLog, formattedLog.Msg)
	//sysLog.Emerg("And this is a daemon emergency with demotag.")

}

func main() {
	//TODO: improve configuration (default priority/facility)
	dest := flag.String("dest", "10.0.0.7:514", "syslog host name")
	protocol := flag.String("protocol", "tcp", "protocol for syslog: tcp or udp")
	port := flag.String("port", ":38090", "port for this webserver")
	flag.Parse()
	config = Config{SyslogHost: *dest, Protocol: *protocol, Port: *port}
	http.HandleFunc("/status", status)
	http.HandleFunc("/", handle)
	log.Println("listen on", *port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
