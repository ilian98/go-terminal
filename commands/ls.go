package commands

import (
	"os"
	"strconv"
	"time"
)

// Ls is a structure for ls command, implementing ExecuteCommand interface
type Ls struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (l *Ls) GetName() string {
	return "ls"
}

// GetPath is a getter for path
func (l *Ls) GetPath() string {
	return l.path
}

// Clone is a method for cloning ls command
func (l *Ls) Clone() ExecuteCommand {
	clone := *l
	return &clone
}

// InitChannel is a method for initializing stopExecution channel
func (l *Ls) InitChannel() {
	l.stopExecution = make(chan struct{}, 1)
}

// StopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (l *Ls) StopSignal() {
	l.stopExecution <- struct{}{}
}

// IsStopSignal is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (l *Ls) IsStopSignal() bool {
	select {
	case <-l.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of ls command
func (l *Ls) Execute(cp CommandProperties) error {
	l.path = cp.Path
	_, outputFile := cp.InputFile, cp.OutputFile

	path, err := os.Open(l.path)
	if err != nil {
		return err
	}
	files, err := path.Readdir(0)
	if err != nil {
		return err
	}
	path.Close()

	lOption := false
	for _, option := range cp.Options {
		if option == "l" {
			lOption = true
		}
	}
	if lOption == false {
		for _, file := range files {
			checkWrite(l, outputFile, file.Name())
			if file.IsDir() {
				checkWrite(l, outputFile, string(os.PathSeparator))
			}
			outputFile.WriteString("    ")
		}
		return nil
	}

	var maxNumberOfDigs int
	for _, file := range files {
		if numberOfDigs := len(strconv.Itoa(int(file.Size()))); maxNumberOfDigs < numberOfDigs {
			maxNumberOfDigs = numberOfDigs
		}
	}
	for _, file := range files {
		checkWrite(l, outputFile, file.Mode().String())
		checkWrite(l, outputFile, " ")

		fileSize := strconv.Itoa(int(file.Size()))
		for i := 0; i < (maxNumberOfDigs - len(fileSize)); i++ {
			checkWrite(l, outputFile, " ")
		}
		checkWrite(l, outputFile, strconv.Itoa(int(file.Size())))

		checkWrite(l, outputFile, " ")
		checkWrite(l, outputFile, outputTime(file.ModTime()))
		checkWrite(l, outputFile, " ")

		checkWrite(l, outputFile, file.Name())
		if file.IsDir() {
			checkWrite(l, outputFile, string(os.PathSeparator))
		}
		checkWrite(l, outputFile, "\n")
	}

	return nil
}

func outputTime(t time.Time) string {
	outputNumber := func(num string) string {
		var output string
		if len(num) == 1 {
			output += "0"
		}
		output += num
		return output
	}

	output := outputNumber(strconv.Itoa(t.Hour())) + ":" + outputNumber(strconv.Itoa(t.Minute())) + " "
	output += outputNumber(strconv.Itoa(t.Day())) + " "
	month := t.Month().String()
	output += month[:3]
	return output
}
