package main

import (
	"errors"
	"strings"
)

type LogEntry struct {
	DateAndTime string
	LogLevel    string
	Message     string
}

func logLineToEntry(line string) (LogEntry, error) {
	var entry LogEntry

	dateLen := len("2023-09-12 17:39:54")

	if len(line) < dateLen {
		return entry, errors.New("No log date")
	}

	entry.DateAndTime = line[:dateLen]

	logLevelStartIdx := strings.Index(line, "[")
	logLevelEndIdx := strings.Index(line, "]")

	if logLevelStartIdx == -1 || logLevelEndIdx == -1 {
		return entry, errors.New("No log level")
	}

	entry.LogLevel = line[logLevelStartIdx+1 : logLevelEndIdx]
	entry.Message = line[logLevelEndIdx+2:]

	return entry, nil
}

func getUsernameAndIP(message string) (string, string, error) {
	if strings.HasPrefix(message, "Disconnecting ") {
		message = message[len("Disconnecting "):]
	}

	nUsernameChars := 0
	for i := range message {
		if string(message[i]) == " " {
			break
		}

		nUsernameChars++
	}

	// Minecraft usernames are at a minimum 3 characters
	if nUsernameChars < 3 {
		return "", "", errors.New("")
	}

	ipFieldStartIdx := strings.Index(message, "[/")
	ipFieldEndIdx := strings.Index(message, "]")

	ipEndIdx := strings.Index(message, ":")

	if ipFieldStartIdx == -1 || ipFieldEndIdx == -1 || ipEndIdx == -1 {
		return "", "", errors.New("")
	}

	ip := message[ipFieldStartIdx+2 : ipEndIdx]

	username, _, found := strings.Cut(message, " ")

	/*if username == "ivattt" {
		fmt.Println(message)
		fmt.Println(username)
	}*/

	if !found {
		return "", "", errors.New("")
	}

	return username, ip, nil
}

// TODO Make it return both username and IP when present
/*func logMessageHasUsernameAndIP(message string) (bool, string) {
	if strings.HasPrefix(message, "Disconnecting ") {
		message = message[len("Disconnecting ")+1:]
	}

	nUsernameChars := 0
	for i := range message {
		if string(message[i]) == " " {
			break
		}

		nUsernameChars++
	}

	// Minecraft usernames are at a minimum 3 characters
	if nUsernameChars < 3 {
		return false, ""
	}

	if !strings.Contains(message, "[/") || !strings.Contains(message, "]") {
		return false, ""
	}

	username, _, found := strings.Cut(message, " ")
	if !found {
		return false, ""
	}

	return true, username
}*/
