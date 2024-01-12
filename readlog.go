package main

import (
	"errors"
	"strings"
	"time"
)

type LogEntry struct {
	DateAndTime string
	LogLevel    string
	Message     string
}

func dateToEpoch(date string) (int, error) {
	time, err := time.Parse(time.DateTime, date)

	if err != nil {
		return 0, err
	}

	return int(time.Unix()), nil
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
	} else if ipFieldStartIdx+2 >= len(message) || ipEndIdx >= len(message) {
		return "", "", errors.New("Invalid log line")
	}

	ip := message[ipFieldStartIdx+2 : ipEndIdx]

	username, _, found := strings.Cut(message, " ")

	if !found {
		return "", "", errors.New("")
	}

	return username, ip, nil
}

// Returns Username, join=0 disconnect=1 disconnect all (server stopped)=2, error
func getUsernameJoinOrLeave(message string) (string, int, error) {
	if message == "Stopping server" || strings.HasPrefix(message, "Server started in ") {
		return "", 2, nil
	}

	ipFieldStartIdx := strings.Index(message, "[/")
	ipFieldEndIdx := strings.Index(message, "]")
	ipEndIdx := strings.Index(message, ":")

	spaceIdx := strings.Index(message, " ")
	if spaceIdx == -1 {
		return "", 0, errors.New("No whitespace in message")
	}

	username := message[:spaceIdx]

	// Minecraft usernames are at a minimum 3 characters
	if len(username) < 3 {
		return "", 0, errors.New("Length of first word in message less than 3")
	}

	if (ipFieldStartIdx != -1) && (ipFieldEndIdx != -1) && (ipEndIdx != -1) {
		if strings.Contains(message, "] logged in with entity id ") {
			return username, 0, nil
		}
	}

	if strings.Contains(message, " lost connection: ") {
		return username, 1, nil
	}

	usernameFromEnd := message[strings.LastIndex(message, " ")+1:]

	if (strings.HasPrefix(message, "CONSOLE: Banning ") && !strings.Contains(message, "CONSOLE: Banning ip")) || strings.HasPrefix(message, "CONSOLE: Kicking ") {
		return usernameFromEnd, 1, nil
	}

	if (strings.Contains(message, ": Banning ") && !strings.Contains(message, ": Banning ip ")) || strings.Contains(message, ": §eKicked player ") {
		if message[len(message)-1] == '!' {
			return usernameFromEnd[:len(usernameFromEnd)-1], 1, nil
		}

		return usernameFromEnd, 1, nil
	}

	return "", 0, errors.New("Not a join/disconnect line")
}
