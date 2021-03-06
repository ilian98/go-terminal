package parser

import (
	"errors"
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

// helper functions for Command struct
func (c1 *Command) equal(c2 Command) bool {
	if c1.Name != c2.Name {
		return false
	}
	if len(c1.Options) != len(c2.Options) {
		return false
	}
	for ind := range c1.Options {
		if c1.Options[ind] != c2.Options[ind] {
			return false
		}
	}
	if len(c1.Arguments) != len(c2.Arguments) {
		return false
	}
	for ind := range c1.Arguments {
		if c1.Arguments[ind] != c2.Arguments[ind] {
			return false
		}
	}

	if c1.Input != c2.Input {
		return false
	}
	if c1.Output != c2.Output {
		return false
	}
	if c1.BgRun != c2.BgRun {
		return false
	}
	return true
}
func (c1 *Command) notEqual(c2 Command) bool {
	return !c1.equal(c2)
}

func commandToString(c Command) string {
	output := c.Name + " ["
	for _, argument := range c.Arguments {
		output += " " + argument
	}
	output += " ] ["
	for _, option := range c.Options {
		output += " " + option
	}
	output += " ]"

	if c.Input != "" {
		output += " " + c.Input
	} else {
		output += " stdin"
	}
	if c.Output != "" {
		output += " " + c.Output
	} else {
		output += " stdout"
	}

	if c.BgRun == true {
		output += " background run"
	}
	return output
}
func newCommand(Name string, Arguments []string, Options []string) Command {
	return Command{Name, Arguments, Options, "", "", false}
}

func testingParseCommandText(t *testing.T, commandText string, expectedResult Command) {
	result, err := parseCommandText(commandText)
	if err != nil {
		t.Errorf("Expected no error, but got: %v\n", err)
		return
	}
	if result.notEqual(expectedResult) {
		t.Errorf("Expected\n")
		t.Error(commandToString(expectedResult))
		t.Errorf("but got\n")
		t.Error(commandToString(result))
		return
	}
}
func TestParseCommandText(t *testing.T) {
	var tests = []struct {
		commandText string
		result      Command
	}{
		{"exit", newCommand("exit", []string{}, []string{})},

		{"ls -l", newCommand("ls", []string{}, []string{"l"})},
		{"ls 	 -l", newCommand("ls", []string{}, []string{"l"})},
		{"ls -l arg1", newCommand("ls", []string{"arg1"}, []string{"l"})},
		{`ls -l arg1 -a "ab c|d"`, newCommand("ls", []string{"arg1", "ab c|d"}, []string{"l", "a"})},
		{`ls -l arg1 "-arg2"`, newCommand("ls", []string{"arg1", "-arg2"}, []string{"l"})},
		{`ls -l arg1 "arg2`, newCommand("ls", []string{"arg1", `"arg2`}, []string{"l"})},
		{`ls -l ""`, newCommand("ls", []string{`""`}, []string{"l"})},

		{"cat <file.txt", Command{"cat", []string{}, []string{}, "file.txt", "", false}},
		{`cat <"file 1.txt"`, Command{"cat", []string{}, []string{}, "file 1.txt", "", false}},
		{`cat >"file 2.txt"`, Command{"cat", []string{}, []string{}, "", "file 2.txt", false}},
		{"cat <file1.txt <file2.txt >file3.txt >file4.txt", Command{"cat", []string{"<file2.txt", ">file4.txt"}, []string{}, "file1.txt", "file3.txt", false}},

		{"ls -l >output.txt &", Command{"ls", []string{}, []string{"l"}, "", "output.txt", true}},
		{"ls -l & >output.txt", Command{"ls", []string{}, []string{"l"}, "", "output.txt", true}},

		{"pwd - < >", Command{"pwd", []string{"-"}, []string{}, ">", "", false}},
		{`pwd - "<" ">"`, Command{"pwd", []string{"-", "<", ">"}, []string{}, "", "", false}},

		{"ls < in.txt > out.txt < in2.txt > out2.txt", Command{"ls", []string{"<", "in2.txt", ">", "out2.txt"}, []string{}, "in.txt", "out.txt", false}},
		{"ls < in.txt >", Command{"ls", []string{}, []string{}, "in.txt", "", false}},
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
func testingParse(t *testing.T, text string, expectedResult []Command, expectedErr error) {
	result, err := Parse(text)
	if err != nil {
		if !errors.Is(err, expectedErr) {
			t.Errorf("Expected %v, but got: %v.\n", expectedErr, err)
			return
		}
		if result != nil {
			t.Errorf("Expected result from parsing to be nil\n")
			return
		}
		return
	}
	for ind := range result {
		if result[ind].notEqual(expectedResult[ind]) {
			t.Errorf("Expected\n")
			t.Error(commandsToString(expectedResult))
			t.Errorf("but got\n")
			t.Error(commandsToString(result))
			return
		}
	}
}
func TestParse(t *testing.T) {
	var tests = []struct {
		text   string
		result []Command
		err    error
	}{
		{"exit", []Command{newCommand("exit", []string{}, []string{})}, nil},
		{`ls -l | cat file1.txt "file 2.txt"`,
			[]Command{
				newCommand("ls", []string{}, []string{"l"}),
				newCommand("cat", []string{"file1.txt", "file 2.txt"}, []string{}),
			},
			nil},
		{`ls -l | cat file1.txt "file 2.txt" >"file 3.txt"`,
			[]Command{
				newCommand("ls", []string{}, []string{"l"}),
				{"cat", []string{"file1.txt", "file 2.txt"}, []string{}, "", "file 3.txt", false},
			},
			nil},
		{`c1 "|" |c2 | c3`,
			[]Command{
				newCommand("c1", []string{"|"}, []string{}),
				newCommand("c2", []string{}, []string{}),
				newCommand("c3", []string{}, []string{}),
			},
			nil},
		{`c1 "|" |c2 & | c3`,
			[]Command{
				newCommand("c1", []string{"|"}, []string{}),
				{"c2", []string{}, []string{}, "", "", true},
				newCommand("c3", []string{}, []string{}),
			},
			nil},
		{"", nil, ErrEmptyCommand},
		{"pwd |   | ls -l ", nil, ErrEmptyCommand},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("Parse(%s)", test.text), func(t *testing.T) {
			testingParse(t, test.text, test.result, test.err)
		})
	}
}

func ExampleParse() {
	parsedCommand, _ := Parse("ls -l &\n")
	fmt.Println(commandsToString(parsedCommand))
	// Output:
	// ls [ ] [ l ] stdin stdout background run
}
