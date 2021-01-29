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
		commands.CloseInputOutputFiles(inputFile, outputFile)
		return ExitCommand
	}
	// check if command with the parsed name exists
	result, ind := i.checkForCommand(i.shellCommandsName, parsedCommand.Name)
	if result == false {
		commands.CloseInputOutputFiles(inputFile, outputFile)
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

// Status is struct used for storing code and command name after InterpretCommand
type Status struct {
	Code    int
	Command string
}

// InterpretCommand is a method of Interpreter that interpretes parsed command and executes command
func (i *Interpreter) InterpretCommand(parsedCommand []parser.Command) []Status {
	type pipe struct {
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
	for _, c := range parsedCommand { // we check if the pipe (command) should be run in bg mode
		if c.BgRun == true {
			bgRun = true
			break
		}
	}
	statuses := make(chan Status, len(parsedCommand))
	origInterpreter := *i // we copy the interpreter to not let path change in pipe
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
			// we make sure that command ran in bg mode won't try to read from stdin
			r, w, err := os.Pipe()
			if err != nil {
				fmt.Printf("%v\n", err)
			} else {
				inputFile = r
			}
			w.Close()
		}

		go func(currInterpreter Interpreter, c parser.Command, inputFile *os.File, outputFile *os.File) {
			statuses <- Status{currInterpreter.ExecuteCommand(c, inputFile, outputFile), c.Name}
			if !isPipe && c.BgRun == false { // path can be changed only for one command not in pipe and bg run
				i.Path = currInterpreter.Path // we don't have concurrent access to i.Path because it isn't pipe
			}
		}(origInterpreter, c, inputFile, outputFile)
	}

	var result []Status
	for i := 0; i < len(parsedCommand); i++ {
		result = append(result, <-statuses)
	}
	return result
}

func (i *Interpreter) checkForCommand(names []string, target string) (bool, int) {
	for i, name := range names {
		if name == target {
			return true, i
		}
	}
	return false, -1
}

func (i *Interpreter) openInputFile(fileName string) (*os.File, error) {
	if fileName == "" {
		return os.Stdin, nil
	}
	file, err := os.Open(commands.FullFileName(i.Path, fileName))
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
	file, err := os.OpenFile(commands.FullFileName(i.Path, fileName), os.O_CREATE|os.O_WRONLY, 0666)
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
