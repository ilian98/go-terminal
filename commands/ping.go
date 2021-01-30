package commands

import (
	"errors"
	"fmt"
	"net"
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
func (p *Ping) Execute(cp *CommandProperties) error {
	p.path = cp.Path
	_, outputFile := cp.InputFile, cp.OutputFile

	if len(cp.Arguments) != 1 {
		return ErrPingOneArg
	}

	p.host = cp.Arguments[0]
	port := "80"

	if err := cp.checkWrite(outputFile, "Pinging "+p.host); err != nil {
		return err
	}

	outputIP := func(connection net.Conn) error {
		address := connection.RemoteAddr().String()
		if err := cp.checkWrite(outputFile, " ["+address+"]\n"); err != nil {
			return err
		}
		return nil
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
					if err := outputIP(result.Conn); err != nil {
						return err
					}
				}
				if err := cp.checkWrite(outputFile, "\n"); err != nil {
					return err
				}
				return fmt.Errorf("%v, %w", result.error, ErrPingDial)
			}
			if p.connection == nil {
				p.connection = result.Conn
			}
		case <-time.After(DefaultTimeOut):
			if i == 0 {
				if err := cp.checkWrite(outputFile, "\n"); err != nil {
					return err
				}

			}
			if err := cp.checkWrite(outputFile, "Request timed out.\n"); err != nil {
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
			if err := cp.checkWrite(outputFile, "Reply from "+p.connection.RemoteAddr().String()+
				": time = "+time.String()+"\n"); err != nil {
				return err
			}
			times = append(times, time)
		}
	}

	if err := p.outputStatistics(cp, len(times), times); err != nil {
		return err
	}
	return nil
}

func (p *Ping) outputStatistics(cp *CommandProperties, cntReceived int, times []time.Duration) error {
	outputFile := cp.OutputFile
	if err := cp.checkWrite(outputFile, "Ping statistics for "); err != nil {
		return err
	}
	if p.connection != nil {
		if err := cp.checkWrite(outputFile, p.connection.RemoteAddr().String()+":\n"); err != nil {
			return err
		}
	} else {
		if err := cp.checkWrite(outputFile, p.host+":\n"); err != nil {
			return err
		}
	}
	lost := PingRepetitions - cntReceived
	if err := cp.checkWrite(outputFile, "    Packets: Sent = "+strconv.Itoa(PingRepetitions)); err != nil {
		return err
	}
	if err := cp.checkWrite(outputFile, ", Received = "+strconv.Itoa(cntReceived)); err != nil {
		return err
	}
	if err := cp.checkWrite(outputFile, ", Lost = "+strconv.Itoa(lost)); err != nil {
		return err
	}
	if err := cp.checkWrite(outputFile, " ("+strconv.Itoa(lost*100/PingRepetitions)+"% loss)\n"); err != nil {
		return err
	}

	if cntReceived == 0 {
		return nil
	}
	if err := cp.checkWrite(outputFile, "Approximate round trip times in milli-seconds:\n"); err != nil {
		return err
	}
	if err := cp.checkWrite(outputFile, "    Minimum = "+minimumTime(times).String()); err != nil {
		return err
	}
	if err := cp.checkWrite(outputFile, ", Maximum = "+maximumTime(times).String()); err != nil {
		return err
	}
	if err := cp.checkWrite(outputFile, ", Average = "+averageTime(times).String()); err != nil {
		return err
	}

	return nil
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
