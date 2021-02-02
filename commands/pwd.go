package commands

// Pwd is a structure for pwd command, implementing ExecuteCommand interface
type Pwd struct {
	path          string
	stopExecution chan struct{}
}

// GetName is a getter for command name
func (p *Pwd) GetName() string {
	return "pwd"
}

// GetPath is a getter for path
func (p *Pwd) GetPath() string {
	return p.path
}

// Clone is a method for cloning pwd command
func (p *Pwd) Clone() ExecuteCommand {
	clone := *p
	return &clone
}

// InitStopCatching is a method for initializing stopExecution channel
func (p *Pwd) InitStopCatching() {
	p.stopExecution = make(chan struct{}, 1)
}

// StopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (p *Pwd) StopSignal() {
	p.stopExecution <- struct{}{}
}

// IsStopSignal is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (p *Pwd) IsStopSignal() bool {
	select {
	case <-p.stopExecution:
		return true
	default:
		return false
	}
}

// Execute is go implementation of pwd command
func (p *Pwd) Execute(cp CommandProperties) error {
	p.path = cp.Path
	_, outputFile := cp.InputFile, cp.OutputFile

	if err := checkWrite(p, outputFile, p.path); err != nil {
		return err
	}

	return nil
}
