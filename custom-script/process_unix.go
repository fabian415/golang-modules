//go:build !windows
// +build !windows

package main

import (
	"os/exec"
	"syscall"
)

// killProcess sends SIGINT to the process on Unix systems
func killProcess(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		return syscall.Kill(cmd.Process.Pid, syscall.SIGINT)
	}
	return nil
}

