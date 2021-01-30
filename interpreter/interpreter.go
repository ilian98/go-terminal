// Package interpreter package interpretes command based on its properties and runs it.
// This is the most important package and can be run almost independentely from the other packages.
// Only InterpretCommand has parameter of struct parser.Command and main struct Interpreter has a field that is slice of interface commands.ExecuteCommand.
//
// Method InterpretCommand should be called with parse information of command stored by package parser.
// Then it makes a potential pipe of commands each of which is executed with method ExecuteCommand.
//
// Method ExecuteCommand only has information about the command that should be executed with any command that is already registered.
// Methods RegisterExitCommand and RegisterCommand are for registering new commands in the interpreter.
package interpreter

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"

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
		return fmt.Errorf("%s - %w", name, ErrCommandExists)
	}
	i.exitCommands = append(i.exitCommands, name)
	return nil
}

// RegisterCommand is a method of Interpreter that can be used to add new command implementing commands.ExecuteCommand interface in shellCommands
func (i *Interpreter) RegisterCommand(c commands.ExecuteCommand) error {
	name := c.GetName()
	if res, _ := i.checkForCommand(i.shellCommandsName, name); res == true {
		return fmt.Errorf("%s - %w", name, ErrCommandExists)
	}
	i.shellCommandsName = append(i.shellCommandsName, name)
	i.shellCommands = append(i.shellCommands, c)
	return nil
}

// These constants are used for the status of method ExecutedCommand
const (
	// CmdInterrupted indicates that the execution of the command was interrupted by Ctrl+C
	CmdInterrupted = -1
	// Ok indicates that there were no errors excluding errors from executed command
	Ok = iota
	// ExitCommand indicates that the command is for exiting the terminal
	ExitCommand
	// InvalidCommandName indicates that the command's parsed name is not present in shellCommandsName
	InvalidCommandName
)

// ExecuteCommand is a method of Interpreter that executes one command given sufficient information after interpreting parsed command.
//
// This method can run the command in background mode if in the parameters bgRun is true.
// Also this method can catch os.Interrupt and alters its default behaviour.
// After the catch, it sends signal to the command that is currently running by writing to its StopExecution channel and then exits the current go routine to call the defer calls closing the opened files!
func (i *Interpreter) ExecuteCommand(name string, arguments []string, options []string, inputFile *os.File, outputFile *os.File, bgRun bool) int {
	// check if command is for exiting the terminal
	if result, _ := i.checkForCommand(i.exitCommands, name); result == true {
		closeInputOutputFiles(inputFile, outputFile)
		return ExitCommand
	}
	// check if command with the parsed name exists
	result, ind := i.checkForCommand(i.shellCommandsName, name)
	if result == false {
		closeInputOutputFiles(inputFile, outputFile)
		return InvalidCommandName
	}

	// Transform command information to element of struct commands.CommandProperties passed by reference to Execute method of command
	cp := commands.CommandProperties{
		Path:       i.Path,
		Arguments:  arguments,
		Options:    options,
		InputFile:  inputFile,
		OutputFile: outputFile,
	}

	command := i.shellCommands[ind].Clone() // we are cloning command so that it runs clean i.e. in initial state
	runCommand := func(cp commands.CommandProperties, bgRun bool) {
		defer closeInputOutputFiles(inputFile, outputFile) // when function ends, then the command stopped and we have to close the opened files

		if bgRun == false { // if we are not in background mode, we should catch Ctrl+C
			result := make(chan error, 1) // we make a channel for waiting result
			signalInterrupt := make(chan os.Signal, 1)
			go func() { // we run the command in new go routine to be able to catch os.Intterupt in current go routine
				result <- command.Execute(cp)
			}()
			signal.Notify(signalInterrupt, os.Interrupt)
			select {
			case <-signalInterrupt:
				command.StopSignal() // we send stop signal to executing command
				runtime.Goexit()     // we stop current goroutine
			case err := <-result:
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
		} else { // in background mode we don't catch Ctrl+C
			err := command.Execute(cp)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}

	if bgRun == true {
		go runCommand(cp, bgRun)
	} else {
		runCommand(cp, bgRun)
		i.Path = command.GetPath() // path changed only when command is not run in background mode
	}

	return Ok
}

// Status is a struct used for storing code and command name after InterpretCommand
type Status struct {
	Code    int
	Command string
}

// InterpretCommand is a method of Interpreter that interpretes parsed command and executes commands.
// It returns a slice with the statuses returned from method ExecuteCommand for every command
//
// If slice parameter length is more than one then a pipe is made.
// All commands are run in background mode if there is at least one command which should be run in background mode.
// Otherwise they are run in normal mode.
func (i *Interpreter) InterpretCommand(parsedCommand []parser.Command) []Status {
	type pipe struct { // structure for grouping read and write end of os.Pipe
		r *os.File
		w *os.File
	}
	var pipes []pipe
	for i := 0; i+1 < len(parsedCommand); i++ {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		pipes = append(pipes, pipe{r, w})
	}

	bgRun := false
	for _, c := range parsedCommand { // we check if there is a command that should be run in background mode
		if c.BgRun == true {
			bgRun = true
			break
		}
	}

	statuses := make(chan Status, len(parsedCommand)) // channel for collecting the statuses of ran commands
	copyInterpreter := *i                             // we copy the interpreter to not let path change in potential pipe
	isPipe := false
	if len(parsedCommand) > 1 {
		isPipe = true
	}
	for ind, c := range parsedCommand {
		inputFile, outputFile, err := i.openInputOutputFiles(c.Input, c.Output)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		if ind > 0 && inputFile == os.Stdin {
			inputFile = pipes[ind-1].r
		} else if ind > 0 {
			pipes[ind-1].r.Close() // closing that end of pipe because it won't be used
		}
		if ind+1 < len(parsedCommand) && outputFile == os.Stdout {
			outputFile = pipes[ind].w
		} else if ind+1 < len(parsedCommand) {
			pipes[ind].w.Close() // closing that end of pipe because it won't be used
		}

		c.BgRun = bgRun
		if c.BgRun == true && inputFile == os.Stdin {
			// we make sure that command ran in background mode won't read from stdin
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				inputFile = r
			}
			w.Close()
		}

		go func(currInterpreter Interpreter, c parser.Command, inputFile *os.File, outputFile *os.File) {
			s := Status{CmdInterrupted, c.Name}
			defer func(s *Status) { // we run this function in defer to write code for command if go routine was exited
				if s.Code == CmdInterrupted { // if code is CmdInterrupted, then the go routine was interrupted
					s.Code = Ok
					statuses <- *s
				}
			}(&s)

			s = Status{currInterpreter.ExecuteCommand(
				c.Name, c.Arguments, c.Options, inputFile, outputFile, c.BgRun,
			), c.Name}
			statuses <- s
			if !isPipe && c.BgRun == false { // path can be changed only for one command not in pipe and bg run
				i.Path = currInterpreter.Path // we don't have concurrent access to i.Path because it isn't pipe
			}
		}(copyInterpreter, c, inputFile, outputFile)
	}

	var result []Status
	for i := 0; i < len(parsedCommand); i++ { // we collect the statuses from the commands
		result = append(result, <-statuses)
	}
	signal.Reset(os.Interrupt) // we remove catching Ctrl+C when all results are collected
	return result
}

// checkForCommand is function for checking if a command name target is present in slice parameter names
func (i *Interpreter) checkForCommand(names []string, target string) (bool, int) {
	for i, name := range names {
		if name == target {
			return true, i
		}
	}
	return false, -1
}
