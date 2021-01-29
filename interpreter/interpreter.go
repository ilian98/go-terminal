package interpreter

import (
	"errors"
	"fmt"
	"os"

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

// ExecuteCommand is a method of Interpreter that executes command given command information after parsing and input and output files which could be stdin and stdout
func (i *Interpreter) ExecuteCommand(parsedCommand parser.Command, inputFile *os.File, outputFile *os.File) int {
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
		Path:       i.Path,
		Arguments:  parsedCommand.Arguments,
		Options:    parsedCommand.Options,
		InputFile:  inputFile,
		OutputFile: outputFile,
	}

	command := i.shellCommands[ind].Clone()
	runCommand := func() {
		if err := command.Execute(cp); err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	if parsedCommand.BgRun == true {
		go runCommand()
	} else {
		runCommand()
		i.Path = command.GetPath() // path changed only when command is not run in bg mode
	}

	return Ok
}

// InterpretCommand is a method of Interpreter that interpretes parsed command and executes command
func (i *Interpreter) InterpretCommand(parsedCommand []parser.Command) int {
	if len(parsedCommand) == 1 {
		inputFile, outputFile, err := i.openInputOutputFiles(parsedCommand[0].Input, parsedCommand[0].Output)
		if parsedCommand[0].BgRun == true && inputFile == os.Stdin {
			// we make sure that command ran in bg mode won't try to read from stdin
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				inputFile = r
			}
			w.Close()
		}

		if err != nil {
			fmt.Printf("%v\n", err)
			return 0
		}
		return i.ExecuteCommand(parsedCommand[0], inputFile, outputFile)
	}
	return 0
}

func (i *Interpreter) checkForCommand(names []string, target string) (bool, int) {
	for i, name := range names {
		if name == target {
			return true, i
		}
	}
	return false, -1
}

func (i *Interpreter) pathFile(fileName string) string {
	return i.Path + string(os.PathSeparator) + fileName
}

func (i *Interpreter) openInputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(i.pathFile(fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("File for reading the input with name %s does not exist", fileName)
		}
		return nil, err
	}
	return file, nil
}

func (i *Interpreter) openOutputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdout, nil
	}
	file, err := os.OpenFile(i.pathFile(fileName), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (i *Interpreter) openInputOutputFiles(input string, output string) (*os.File, *os.File, error) {
	inputFile, err := i.openInputFile(input)
	if err != nil {
		return nil, nil, err
	}

	outputFile, err := i.openOutputFile(output)
	if err != nil {
		return nil, nil, err
	}

	return inputFile, outputFile, nil
}
