package parser

import (
	"fmt"
	"testing"
)

func testingReplaceEnclosed(t *testing.T, text string, target byte, value byte, expectedResult string) {
	if result := replaceEnclosed(text, target, value); result != expectedResult {
		t.Errorf("Expected %s but got %s", expectedResult, result)
	}
}
func TestReplaceEnclosed(t *testing.T) {
	var tests = []struct {
		text          string
		target, value byte
		result        string
	}{
		{"abcd abcd", '|', 0, "abcd abcd"},
		{"abcd|abcd", '|', 0, "abcd|abcd"},
		{`abcd "a d jk" abcd`, ' ', 0, `abcd "` + "a\x00d\x00jk" + `" abcd`},
		{`abcd "a b c" "d e" abcd`, ' ', 0, `abcd "a` + "\x00b\x00c" + `" "d` + "\x00e" + `" abcd`},
		{`abcd "` + "a\x00d\x00jk" + `" abcd`, 0, ' ', `abcd "a d jk" abcd`},
		{"", '|', 0, ""},
		{`abcd "abc|bc" | ab`, '|', '+', `abcd "abc+bc" | ab`},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("replaceEnclosed(%s,%c,%c)", test.text, test.target, test.value), func(t *testing.T) {
			testingReplaceEnclosed(t, test.text, test.target, test.value, test.result)
		})
	}
}

func commandToString(c Command) string {
	output := c.Name + " ["
	for _, option := range c.Options {
		output += " " + option
	}
	output += " ] ["
	for _, argument := range c.Arguments {
		output += " " + argument
	}
	output += " ]"
	return output
}
func testingParseCommandText(t *testing.T, commandText string, expectedResult Command) {
	if result := parseCommandText(commandText); result.NotEqual(expectedResult) {
		t.Errorf("Expected\n")
		t.Error(commandToString(expectedResult))
		t.Errorf("but got\n")
		t.Error(commandToString(result))
	}
}
func TestParseCommandText(t *testing.T) {
	var tests = []struct {
		commandText string
		result      Command
	}{
		{"exit", Command{"exit", []string{}, []string{}}},
		{"ls -l", Command{"ls", []string{"l"}, []string{}}},
		{"ls -l arg1", Command{"ls", []string{"l"}, []string{"arg1"}}},
		{"ls -l arg1 -a " + `"ab c|d"`, Command{"ls", []string{"l", "a"}, []string{"arg1", "ab c|d"}}},
		{"ls -l arg1 " + `"-arg2"`, Command{"ls", []string{"l"}, []string{"arg1", "-arg2"}}},
		{"ls -l arg1 " + `"arg2`, Command{"ls", []string{"l"}, []string{"arg1", `"arg2`}}},
		{"ls -l " + `""`, Command{"ls", []string{"l"}, []string{`""`}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("parseCommandText(%s)", test.commandText), func(t *testing.T) {
			testingParseCommandText(t, test.commandText, test.result)
		})
	}
}

func commandsToString(commands []Command) string {
	output := commandToString(commands[0])
	for _, command := range commands[1:] {
		output += " | " + commandToString(command)
	}
	return output
}
func testingParse(t *testing.T, text string, expectedResult []Command) {
	result := Parse(text)
	for ind := range result {
		if result[ind].NotEqual(expectedResult[ind]) {
			t.Errorf("Expected\n")
			t.Error(commandsToString(expectedResult))
			t.Errorf("but got\n")
			t.Error(commandsToString(result))
		}
	}
}
func TestParse(t *testing.T) {
	var tests = []struct {
		text   string
		result []Command
	}{
		{"exit", []Command{{"exit", []string{}, []string{}}}},
		{"ls -l | cat file1.txt " + `"file 2.txt"`, []Command{
			{"ls", []string{"l"}, []string{}},
			{"cat", []string{}, []string{"file1.txt", "file 2.txt"}}}},
		{"c1 " + `"|"` + "|c2 | c3", []Command{
			{"c1", []string{}, []string{"|"}},
			{"c2", []string{}, []string{}},
			{"c3", []string{}, []string{}}}},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Parse(%s)", test.text), func(t *testing.T) {
			testingParse(t, test.text, test.result)
		})
	}
}

func ExampleParse() {
	parsedCommand := Parse("ls -l\n")
	fmt.Println(commandsToString(parsedCommand))
	// Output:
	// ls [ l ] [ ]
}
