package terminator

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	// Timeout defines the amount of time to wait before exiting the application after a termination request
	Timeout = 30 * time.Second

	// WaitGroup is the global semaphore used to indicate that a goroutine is running
	waitGroup sync.WaitGroup

	// stopChannel is the global channel used to indicate an application termination request.
	// The channel is closed on the termination event.
	stopChannel chan struct{}

	Signal <-chan struct{}

	// Aliases of waitGroup functions
	Add  = waitGroup.Add
	Done = waitGroup.Done
	Wait = waitGroup.Wait

	// osExit is a variable for testing purposes
	osExit = os.Exit
)

// Start calls a concurrent routine to respond to termination requests.
// In case of an event, the signal channel is closed and after the specified time,
// the application is terminated with exit code 143.
func Start() {
	stopChannel = make(chan struct{})
	Signal = stopChannel

	// Listen to interrupt and termination signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	// Guarantee termination after specified timeout or cancel goroutine if stopped otherwise
	go func() {
		select {
		case <-signalCh:
			Stop()
			time.Sleep(Timeout)
			osExit(143)
		case <-stopChannel:
		}
	}()
}

// Stop closes Signal channel without terminating the application itself.
func Stop() {
	close(stopChannel)
}

// IsStopped can be used to check the current termination state
// if continuous monitoring of the Signal channel is not possible.
func IsStopped() bool {
	select {
	case <-stopChannel:
		return true
	default:
	}
	return false
}
