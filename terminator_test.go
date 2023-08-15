package terminator

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"syscall"
	"testing"
	"time"
)

func prepare(osExitCode *int) {
	Timeout = time.Second
	osExit = func(code int) { *osExitCode = code }
	Start()
}

func TestTerminator_SIGTERM(t *testing.T) {
	osExitCode := -1
	signalClosed := false
	prepare(&osExitCode)

	go func() {
		<-Signal
		signalClosed = true
	}()

	go func() {
		<-time.After(time.Hour)
	}()

	err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	require.NoError(t, err)

	time.Sleep(Timeout + time.Second)
	assert.Equalf(t, 143, osExitCode, "expected '143' exit code instead of %d", osExitCode)
	assert.Equal(t, true, signalClosed, "expected a closed signal channel")
}

func TestTerminator_Wait(t *testing.T) {
	osExitCode := -1
	prepare(&osExitCode)

	Add(1)
	go func() {
		<-Signal
		Done()
	}()

	go func() {
		time.Sleep(time.Second)
		err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		require.NoError(t, err)
	}()

	Wait()
	time.Sleep(Timeout + 2*time.Second)
	assert.Equalf(t, 143, osExitCode, "expected '143' exit code instead of %d", osExitCode)
}

func TestTerminator_IsStopped(t *testing.T) {
	osExitCode := -1
	prepare(&osExitCode)

	assert.Equal(t, false, IsStopped(), "expected an open channel")
	Stop()

	assert.Equal(t, true, IsStopped(), "expected an closed channel")
}
