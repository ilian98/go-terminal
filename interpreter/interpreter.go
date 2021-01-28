package interpreter

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/ilian98/go-terminal/commands"
	"github.com/ilian98/go-terminal/parser"
)

// Interpreter is struct for working with parsed commands
type Interpreter struct {
	Path              string
	exitCommands      []string
	shellCommandsName []string
	shellCommands     []commands.ExecuteCommand
}

var (
	// ErrCommandExists indicates that command is already registered and won't be added
	ErrCommandExists = errors.New("command with that name already exists")
)

// RegisterExitCommand is a method of Interpreter that can be used to add new name in exitCommands
func (i *Interpreter) RegisterExitCommand(name string) error {
	if res, _ := i.checkForCommand(i.exitCommands, name); res == true {
		return ErrCommandExists
	}
	i.exitCommands = append(i.exitCommands, name)
	return nil
}

// RegisterCommand is a method of Interpreter that can be used to add new command implementing commands.ExecuteCommand interface in shellCommands
func (i *Interpreter) RegisterCommand(c commands.ExecuteCommand) error {
	name := c.GetName()
	if res, _ := i.checkForCommand(i.shellCommandsName, name); res == true {
		return ErrCommandExists
	}
	i.shellCommandsName = append(i.shellCommandsName, name)
	i.shellCommands = append(i.shellCommands, c)
	return nil
}

// These constants are used for the status of method ExecutedCommand
const (
	// Ok indicates that there were no errors excluding errors from executed command
	Ok = iota
	// ExitCommand indicates that the command is for exiting the terminal
	ExitCommand
	// InvalidCommandName indicates that the command's parsed name is not present in shellCommandsName
	InvalidCommandName
)

// ExecuteCommand is a method of Interpreter that executes command given command information after parsing
func (i *Interpreter) ExecuteCommand(parsedCommand parser.Command) int {
	// check if command is for exiting the terminal
	if result, _ := i.checkForCommand(i.exitCommands, parsedCommand.Name); result == true {
		return ExitCommand
	}
	// check if command with the parsed name exists
	result, ind := i.checkForCommand(i.shellCommandsName, parsedCommand.Name)
	if result == false {
		return InvalidCommandName
	}

	cp := commands.CommandProperties{
		Path:      i.Path,
		Arguments: parsedCommand.Arguments,
		Options:   parsedCommand.Options,
		Input:     parsedCommand.Input,
		Output:    parsedCommand.Output,
	}

	removeFileName := ""
	if parsedCommand.BgRun == true && cp.Input == "" {
		// we make sure that command ran in bg mode won't try to read from stdin
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomNumber := strconv.Itoa(r.Int())
		removeFileName = cp.Path + string(os.PathSeparator) + "dummy-file-mock-stdin-bg-run" + randomNumber
		dummy, _ := os.Create(removeFileName)
		dummy.Close()

		cp.Input = "dummy-file-mock-stdin-bg-run" + randomNumber
	}

	command := i.shellCommands[ind].Clone()
	runCommand := func(removeFileName string) {
		if err := command.Execute(cp); err != nil {
			fmt.Printf("%v\n", err)
		}
		if removeFileName != "" {
			if err := os.Remove(removeFileName); err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}

	if parsedCommand.BgRun == true {
		go runCommand(removeFileName)
	} else {
		runCommand("")
		i.Path = command.GetPath() // path changed only when command is not run in bg mode
	}

	return Ok
}

func (i *Interpreter) checkForCommand(names []string, target string) (bool, int) {
	for i, name := range names {
		if name == target {
			return true, i
		}
	}
	return false, -1
}
