package commands

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	// ErrPingOneArg indicates that the ping command has different argument count from 1
	ErrPingOneArg = errors.New("One argument is needed")
	// ErrPingDial indicates that dial returned error and suggests possible reasons
	ErrPingDial = errors.New("possible reasons are Internet connectivity problem or host unreachability")
)

// Ping is a structure for ping command, implementing ExecuteCommand interface
type Ping struct {
	path          string
	stopExecution chan struct{}
	host          string
	connection    net.Conn
}

// GetName is a getter for command name
func (p *Ping) GetName() string {
	return "ping"
}

// GetPath is a getter for path
func (p *Ping) GetPath() string {
	return p.path
}

// Clone is a method for cloning ping command
func (p *Ping) Clone() ExecuteCommand {
	clone := *p
	return &clone
}

// InitStopSignalCatching is a method for initializing stopExecution channel
func (p *Ping) InitStopSignalCatching() {
	p.stopExecution = make(chan struct{}, 1)
}

// SendStopSignal is a method for registering stop signal of the execution of the command
// It writes to stopExecution channel
func (p *Ping) SendStopSignal() {
	p.stopExecution <- struct{}{}
}

// IsStopSignalReceived is a method for checking if stop signal was sent
// It checks if there is a signal in stopExecution channel
func (p *Ping) IsStopSignalReceived() bool {
	select {
	case <-p.stopExecution:
		return true
	default:
		return false
	}
}

const (
	// PingRepetitions stores the number of pings that are made
	PingRepetitions = 4
	// DefaultDialTimeOut stores the timeout for dial
	DefaultDialTimeOut = time.Duration(6 * time.Second)
	// DefaultTimeOut stores the timeout for each ping
	DefaultTimeOut = time.Duration(5 * time.Second)
)

// Execute is go implementation of ping command
func (p *Ping) Execute(cp CommandProperties) error {
	p.path = cp.Path
	_, outputFile := cp.InputFile, cp.OutputFile

	if len(cp.Arguments) != 1 {
		return ErrPingOneArg
	}

	p.host = cp.Arguments[0]
	port := "80"

	if err := checkWrite(p, outputFile, "Pinging "+p.host); err != nil {
		return err
	}

	outputIP := func(connection net.Conn) error {
		address := connection.RemoteAddr().String()
		if err := checkWrite(p, outputFile, " ["+address+"]\n"); err != nil {
			return err
		}
		return nil
	}

	var times []time.Duration
	for i := 0; i < PingRepetitions; i++ {
		ch := make(chan struct {
			net.Conn
			error
		}, 1) // channel for collecting possible error
		startTime := time.Now()
		go func() {
			connection, err := net.DialTimeout("tcp", p.host+":"+port, DefaultDialTimeOut)
			ch <- struct {
				net.Conn
				error
			}{connection, err}
		}()
		endTime := startTime // default value for endTime

		select {
		case result := <-ch:
			endTime = time.Now()
			if result.error != nil {
				if i == 0 && result.Conn != nil {
					if err := outputIP(result.Conn); err != nil {
						return err
					}
				}
				if err := checkWrite(p, outputFile, "\n"); err != nil {
					return err
				}
				return fmt.Errorf("%v, %w", result.error, ErrPingDial)
			}
			if p.connection == nil {
				p.connection = result.Conn
			}
		case <-time.After(DefaultTimeOut):
			if i == 0 {
				if err := checkWrite(p, outputFile, "\n"); err != nil {
					return err
				}

			}
			if err := checkWrite(p, outputFile, "Request timed out.\n"); err != nil {
				return err
			}
		}

		if endTime != startTime {
			if i == 0 {
				if err := outputIP(p.connection); err != nil {
					return err
				}
			}
			time := endTime.Sub(startTime)
			if err := checkWrite(p, outputFile, "Reply from "+p.connection.RemoteAddr().String()+
				": time = "+time.String()+"\n"); err != nil {
				return err
			}
			times = append(times, time)
		}
	}

	if err := p.outputStatistics(outputFile, len(times), times); err != nil {
		return err
	}
	return nil
}

// outputStatistics function is helper for writing number of sent, received, lost packets, minTime, maxTime and averageTime of pings
func (p *Ping) outputStatistics(outputFile *os.File, cntReceived int, times []time.Duration) error {
	if err := checkWrite(p, outputFile, "Ping statistics for "); err != nil {
		return err
	}
	if p.connection != nil {
		if err := checkWrite(p, outputFile, p.connection.RemoteAddr().String()+":\n"); err != nil {
			return err
		}
	} else {
		if err := checkWrite(p, outputFile, p.host+":\n"); err != nil {
			return err
		}
	}
	lost := PingRepetitions - cntReceived
	if err := checkWrite(p, outputFile, "    Packets: Sent = "+strconv.Itoa(PingRepetitions)); err != nil {
		return err
	}
	if err := checkWrite(p, outputFile, ", Received = "+strconv.Itoa(cntReceived)); err != nil {
		return err
	}
	if err := checkWrite(p, outputFile, ", Lost = "+strconv.Itoa(lost)); err != nil {
		return err
	}
	if err := checkWrite(p, outputFile, " ("+strconv.Itoa(lost*100/PingRepetitions)+"% loss)\n"); err != nil {
		return err
	}

	if cntReceived == 0 {
		return nil
	}
	if err := checkWrite(p, outputFile, "Approximate round trip times in milli-seconds:\n"); err != nil {
		return err
	}
	if err := checkWrite(p, outputFile, "    Minimum = "+minimumTime(times).String()); err != nil {
		return err
	}
	if err := checkWrite(p, outputFile, ", Maximum = "+maximumTime(times).String()); err != nil {
		return err
	}
	if err := checkWrite(p, outputFile, ", Average = "+averageTime(times).String()); err != nil {
		return err
	}

	return nil
}

// minimumTime function is helper for finding minimum time in the slice parameter
func minimumTime(times []time.Duration) time.Duration {
	min := times[0]
	for _, time := range times[1:] {
		if min > time {
			min = time
		}
	}
	return min
}

// maximumTime function is helper for finding maximum time in the slice parameter
func maximumTime(times []time.Duration) time.Duration {
	max := times[0]
	for _, time := range times[1:] {
		if max < time {
			max = time
		}
	}
	return max
}

// averageTime function is helper for finding average time for the slice parameter
func averageTime(times []time.Duration) time.Duration {
	sum := times[0]
	for _, time := range times[1:] {
		sum += time
	}
	return time.Duration(int(sum.Milliseconds())/len(times)) * time.Millisecond
}
