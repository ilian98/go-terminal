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
	// ErrPingDial indicates that dial return error and suggests possible reasons
	ErrPingDial = errors.New("possible reasons are Internet connectivity problem or host unreachability")
)

// Ping is a structure for ping command, implementing ExecuteCommand interface
type Ping struct {
	path       string
	host       string
	connection net.Conn
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
	inputFile, outputFile := cp.InputFile, cp.OutputFile
	defer closeInputOutputFiles(inputFile, outputFile)

	if len(cp.Arguments) != 1 {
		return ErrPingOneArg
	}

	p.host = cp.Arguments[0]
	port := "80"

	outputFile.WriteString("Pinging " + p.host)

	outputIP := func(connection net.Conn) {
		address := connection.RemoteAddr().String()
		outputFile.WriteString(" [" + address + "]\n")
	}

	var times []time.Duration
	for i := 0; i < PingRepetitions; i++ {
		ch := make(chan struct {
			net.Conn
			error
		}, 1)
		startTime := time.Now()
		go func() {
			connection, err := net.DialTimeout("tcp", p.host+":"+port, DefaultDialTimeOut)
			ch <- struct {
				net.Conn
				error
			}{connection, err}
		}()
		endTime := startTime

		select {
		case result := <-ch:
			endTime = time.Now()
			if result.error != nil {
				if i == 0 && result.Conn != nil {
					outputIP(result.Conn)
				}
				outputFile.WriteString("\n")
				return fmt.Errorf("%v, %w", result.error, ErrPingDial)
			}
			if p.connection == nil {
				p.connection = result.Conn
			}
		case <-time.After(DefaultTimeOut):
			if i == 0 {
				outputFile.WriteString("\n")
			}
			outputFile.WriteString("Request timed out.\n")
		}

		if endTime != startTime {
			if i == 0 {
				outputIP(p.connection)
			}
			time := endTime.Sub(startTime)
			outputFile.WriteString("Reply from " + p.connection.RemoteAddr().String() + ": time = " + time.String() + "\n")
			times = append(times, time)
		}
	}

	p.outputStatistics(outputFile, len(times), times)
	return nil
}

func (p *Ping) outputStatistics(outputFile *os.File, cntReceived int, times []time.Duration) {
	outputFile.WriteString("Ping statistics for ")
	if p.connection != nil {
		outputFile.WriteString(p.connection.RemoteAddr().String() + ":\n")
	} else {
		outputFile.WriteString(p.host + ":\n")
	}
	lost := PingRepetitions - cntReceived
	outputFile.WriteString("    Packets: Sent = " + strconv.Itoa(PingRepetitions))
	outputFile.WriteString(", Received = " + strconv.Itoa(cntReceived))
	outputFile.WriteString(", Lost = " + strconv.Itoa(lost))
	outputFile.WriteString(" (" + strconv.Itoa(lost*100/PingRepetitions) + "% loss)\n")

	if cntReceived == 0 {
		return
	}
	outputFile.WriteString("Approximate round trip times in milli-seconds:\n")
	outputFile.WriteString("    Minimum = " + minimumTime(times).String())
	outputFile.WriteString(", Maximum = " + maximumTime(times).String())
	outputFile.WriteString(", Average = " + averageTime(times).String())
}

func minimumTime(times []time.Duration) time.Duration {
	min := times[0]
	for _, time := range times[1:] {
		if min > time {
			min = time
		}
	}
	return min
}
func maximumTime(times []time.Duration) time.Duration {
	max := times[0]
	for _, time := range times[1:] {
		if max < time {
			max = time
		}
	}
	return max
}
func averageTime(times []time.Duration) time.Duration {
	sum := times[0]
	for _, time := range times[1:] {
		sum += time
	}
	return time.Duration(int(sum.Milliseconds())/len(times)) * time.Millisecond
}
