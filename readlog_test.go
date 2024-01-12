package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetUsernameJoinOrLeave(t *testing.T) {
	type test struct {
		input string

		username string
		action   int
		err      error
	}

	tests := []test{
		{input: "", username: "", action: 0, err: errors.New("No whitespace in message")},

		{input: "Stopping server", username: "", action: 2, err: nil},
		{input: "Server started in 3650 Milliseconds (3 seconds)! For help, type \"help\" or \"?\"", username: "", action: 2, err: nil},

		{input: "kivattt [/69.69.69.69:69696] logged in with entity id 145 at (-42.65625, 89.0, 13.46875)", username: "kivattt", action: 0, err: nil},

		{input: "kivattt lost connection: disconnect.quitting", username: "kivattt", action: 1, err: nil},
		{input: "CONSOLE: Banning kivattt", username: "kivattt", action: 1, err: nil},
		{input: "CONSOLE: Kicking kivattt", username: "kivattt", action: 1, err: nil},
		{input: "kivattt: §eKicked player test!", username: "test", action: 1, err: nil},
		{input: "kivattt: Banning test", username: "test", action: 1, err: nil},

		{input: "Disconnecting testing [/69.69.69.69:69696]: You are banned from this server!", username: "", action: 0, err: errors.New("Not a join/disconnect line")},
	}

	for _, tc := range tests {
		username, action, err := getUsernameJoinOrLeave(tc.input)
		if !reflect.DeepEqual(tc.username, username) {
			t.Fatalf("Expected username: %v, got: %v", tc.username, username)
		}

		if !reflect.DeepEqual(tc.action, action) {
			t.Fatalf("Expected action: %v, got: %v", tc.action, action)
		}

		if !reflect.DeepEqual(tc.err, err) {
			t.Fatalf("Expected error: %v, got: %v", tc.err, err)
		}
	}
}

func TestGetUsernameAndIP(t *testing.T) {
	_, _, err := getUsernameAndIP("2025-11-10 12:10:55 [INFO] CONSOLE: Banning ip [15:19:34] [INFO   ] Disconnecting username [/10.10.10.10:12345]: You are banned fr")

	if err == nil {
		t.Fatalf("Expected error, but got nil")
	}
}
