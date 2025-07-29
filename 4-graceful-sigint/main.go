//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}
	sigkillChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigkillChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the process (blocking)
	go func ()  {
		proc.Run()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-sigkillChan:
		stopDone := make(chan bool, 1)

		go func ()  {
			proc.Stop()
			close(stopDone)
		}()

		select {
		case <-stopDone:
			return
		case <-sigkillChan:
			return
		}
	}
}
