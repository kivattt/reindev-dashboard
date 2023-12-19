package main

import (
	"bufio"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

var templates = template.Must(template.ParseFiles("pages/index.html", "pages/username-to-ips.html", "pages/ip-to-usernames.html", "pages/find-alts.html"))
var serverLogPath = "../server.log"

var nHistoricalUsersCached = 0
var firstLogDate = ""

type FrontPage struct {
	NHistoricalUsers int
	FirstLogDate     string
}

type UsernameToIPsPage struct {
	NoParams bool
	Username string
	IPs      []string
}

type IPToUsernamesPage struct {
	NoParams  bool
	IP        string
	Usernames []string
}

type FindAltsPage struct {
	NoParams bool
	Username string
	AltAccounts []string
}

func getFirstLogDate() string {
	file, err := os.Open(serverLogPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		entry, err := logLineToEntry(line)
		if err != nil {
			continue
		}

		return entry.DateAndTime
	}

	return ""
}

func usernameToIPs(username string) []string {
	if len(username) == 0 || len(username) > 64 {
		return []string{}
	}

	file, err := os.Open(serverLogPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var ips []string

	for scanner.Scan() {
		line := scanner.Text()

		entry, err := logLineToEntry(line)
		if err != nil {
			continue
		}

		if entry.LogLevel != "INFO" {
			continue
		}

		if strings.HasPrefix(entry.Message, username+" [/") || strings.HasPrefix(entry.Message, "Disconnecting "+username+" [/") {
			ipFieldStartIdx := strings.Index(entry.Message, "[/")
			ipFieldEndIdx := strings.Index(entry.Message, "]")

			ipEndIdx := strings.Index(entry.Message, ":")

			if ipFieldStartIdx == -1 || ipFieldEndIdx == -1 || ipEndIdx == -1 {
				continue
			}

			ip := entry.Message[ipFieldStartIdx+2 : ipEndIdx]
			if !slices.Contains(ips, ip) {
				ips = append(ips, ip)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return []string{}
	}

	return ips
}

func ipToUsernames(ip string) []string {
	if len(ip) < 7 || len(ip) > 15 {
		return []string{}
	}

	file, err := os.Open(serverLogPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var usernames []string

	for scanner.Scan() {
		line := scanner.Text()

		entry, err := logLineToEntry(line)
		if err != nil {
			continue
		}

		if entry.LogLevel != "INFO" {
			continue
		}

		username, ipFound, err := getUsernameAndIP(entry.Message)
		if err != nil {
			continue
		}

		if ipFound == ip {
			if !slices.Contains(usernames, username) {
				usernames = append(usernames, username)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return []string{}
	}

	return usernames
}

func getAllHistoricalUsernames() []string {
	file, err := os.Open(serverLogPath)
	if err != nil {
		return []string{}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var usernames []string

	for scanner.Scan() {
		line := scanner.Text()

		entry, err := logLineToEntry(line)
		if err != nil {
			continue
		}

		if entry.LogLevel != "INFO" {
			continue
		}

		username, _, err := getUsernameAndIP(entry.Message)
		if err != nil {
			continue
		}

		if !slices.Contains(usernames, username) {
			usernames = append(usernames, username)
		}
	}

	if err := scanner.Err(); err != nil {
		return []string{}
	}

	nHistoricalUsersCached = len(usernames)

	return usernames
}

func findAlts(targetUsername string) []string {
	ips := usernameToIPs(targetUsername)

	var altAccounts []string

	for _, ip := range ips {
		usernames := ipToUsernames(ip)

		for _, username := range usernames {
			if username == targetUsername {
				continue
			}

			if !slices.Contains(altAccounts, username) {
				altAccounts = append(altAccounts, username)
			}
		}
	}

	return altAccounts
}

func handler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		//		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Access\"")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//		w.Write([]byte("bruh"))
		return
	}

	if username != "kiva" || password != "password" {
		//		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Access\"")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//		w.Write([]byte("bruh"))
		return
	}

	if r.URL.Path == "/" {
		templates.ExecuteTemplate(w, "index.html", FrontPage{NHistoricalUsers: nHistoricalUsersCached, FirstLogDate: firstLogDate})
		return
	}

	if strings.HasPrefix(r.URL.Path, "/username-to-ips") {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		username := strings.Trim(params.Get("username"), " ")

		if len(username) > 64 {
			http.Error(w, "Enter a username less than 64 characters", http.StatusInternalServerError)
			return
		}

		templates.ExecuteTemplate(w, "username-to-ips.html", UsernameToIPsPage{NoParams: username == "", Username: username, IPs: usernameToIPs(username)})
		return
	}

	if strings.HasPrefix(r.URL.Path, "/ip-to-usernames") {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		ipAddress := strings.Trim(params.Get("ipaddress"), " ")

		templates.ExecuteTemplate(w, "ip-to-usernames.html", IPToUsernamesPage{NoParams: ipAddress == "", IP: ipAddress, Usernames: ipToUsernames(ipAddress)})
		return
	}

	if strings.HasPrefix(r.URL.Path, "/find-alts") {
		params, _ := url.ParseQuery(r.URL.RawQuery)
		username := strings.Trim(params.Get("username"), " ")

		templates.ExecuteTemplate(w, "find-alts.html", FindAltsPage{NoParams: username == "", Username: username, AltAccounts: findAlts(username)})
		return
	}

	http.NotFound(w, r)
}

func mainCSSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/main.css")
}

func iconPNGHandler(w http.ResponseWriter, r *http.Request) {
	backgroundImage := "pages/img/good.png"

	rnd := rand.Float64()

	if rnd < 0.05 {
		backgroundImage = "pages/img/bad.png"
	}

	http.ServeFile(w, r, backgroundImage)
}

func bgPNGHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/img/bg.png")
}

func angelPNGHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/img/angel.png")
}

func main() {
	getAllHistoricalUsernames() // Updates the nHistoricalUsersCached
	firstLogDate = getFirstLogDate()

	rand.Seed(time.Now().Unix())

	http.HandleFunc("/", handler)
	http.HandleFunc("/main.css", mainCSSHandler)

	http.HandleFunc("/img/fen.png", iconPNGHandler)
	http.HandleFunc("/img/bg.png", bgPNGHandler)
	http.HandleFunc("/img/angel.png", angelPNGHandler)

	http.ListenAndServe(":8080", nil)
}
