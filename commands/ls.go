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

// InitStopSignalCatching is a method for initializing stopExecution channel
func (l *Ls) InitStopSignalCatching() {
	l.stopExecution = make(chan struct{}, 1)
}

// SendStopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (l *Ls) SendStopSignal() {
	l.stopExecution <- struct{}{}
}

// IsStopSignalReceived is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (l *Ls) IsStopSignalReceived() bool {
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
			if err := checkWrite(l, outputFile, file.Name()); err != nil {
				return err
			}
			if file.IsDir() {
				if err := checkWrite(l, outputFile, string(os.PathSeparator)); err != nil {
					return err
				}
			}
			if err := checkWrite(l, outputFile, "    "); err != nil {
				return err
			}
		}
		return nil
	}

	var maxNumberOfDigs int // we count the maximum number of digits so column with file size will be "aligned" right
	for _, file := range files {
		if numberOfDigs := len(strconv.Itoa(int(file.Size()))); maxNumberOfDigs < numberOfDigs {
			maxNumberOfDigs = numberOfDigs
		}
	}
	for _, file := range files {
		if err := checkWrite(l, outputFile, file.Mode().String()); err != nil { // we write file mode
			return err
		}
		if err := checkWrite(l, outputFile, " "); err != nil {
			return err
		}

		fileSize := strconv.Itoa(int(file.Size()))
		for i := 0; i < (maxNumberOfDigs - len(fileSize)); i++ {
			if err := checkWrite(l, outputFile, " "); err != nil {
				return err
			}
		}
		if err := checkWrite(l, outputFile, strconv.Itoa(int(file.Size()))); err != nil { // we write file size
			return err
		}

		if err := checkWrite(l, outputFile, " "); err != nil {
			return err
		}
		if err := checkWrite(l, outputFile, outputTime(file.ModTime())); err != nil { // we write the data and time of last modification
			return err
		}
		if err := checkWrite(l, outputFile, " "); err != nil {
			return err
		}

		if err := checkWrite(l, outputFile, file.Name()); err != nil { // lastly in row we write file name
			return err
		}
		if file.IsDir() {
			if err := checkWrite(l, outputFile, string(os.PathSeparator)); err != nil {
				return err
			}
		}
		if err := checkWrite(l, outputFile, "\n"); err != nil {
			return err
		}
	}

	return nil
}

// outputTime function is helper for writing time and date in format hh:mm dd mmm
func outputTime(t time.Time) string {
	outputNumber := func(num string) string { // another helper function for writing one-digit number with leading zero
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
