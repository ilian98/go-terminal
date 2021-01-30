package commands

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func testingPing(t *testing.T, arguments []string, expectedResult string, expectedErr error) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal("Fatal error - cannot make pipe! - %w", err)
	}

	ping := Ping{}
	errPing := ping.Execute(&CommandProperties{"", arguments, []string{""}, os.Stdin, w, make(chan struct{}, 1)})

	takeResult := func() string {
		output := make([]byte, 1<<10)
		if _, err := r.Read(output); err != nil {
			t.Fatal("Fatal error - cannot read from pipe! - %w", err)
		}
		text := string(output)
		return text
	}

	if errPing != nil {
		if expectedErr == nil {
			t.Errorf("Expecting no error from Ping function, but got: %w\n", errPing)
			return
		} else if !errors.Is(errPing, expectedErr) {
			t.Errorf("Expecting error %v, bug got: %w", expectedErr, errPing)
			return
		} else if expectedResult != "" {
			text := takeResult()
			if split := strings.SplitN(text, "\n", 2); len(split) > 1 {
				text = split[0]
			}
			if text != expectedResult {
				t.Errorf("Expecting %s, but got: %s", text, expectedResult)
				return
			}
		}
		return
	}

	text := takeResult()
	text = strings.TrimRight(text, "\u0000")
	if len(text) < len(expectedResult) || text[:len(expectedResult)] != expectedResult {
		t.Errorf("Expecting %s, but got: %s", text, expectedResult)
		return
	}
}

func TestPing(t *testing.T) {
	var tests = []struct {
		arguments []string
		result    string
		err       error
	}{
		{[]string{}, "", ErrPingOneArg},
		{[]string{"google.com", "google.com"}, "", ErrPingOneArg},
		{[]string{"1"}, "Pinging 1", ErrPingDial},
		{[]string{"noibg.com"},
			`Pinging noibg.com
Request timed out.
Request timed out.
Request timed out.
Request timed out.
Ping statistics for noibg.com:` + "\n    Packets: Sent = 4, Received = 0, Lost = 4 (100% loss)", nil},
		{[]string{"google.com"}, "Pinging google.com [", nil},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Cat test with arguments %v ", test.arguments), func(t *testing.T) {
			testingPing(t, test.arguments, test.result, test.err)
		})
	}
}
