//go:build windows
// +build windows

package main

import (
	"os/exec"
)

// killProcess terminates the process on Windows
func killProcess(cmd *exec.Cmd) error {
	if cmd != nil && cmd.Process != nil {
		return cmd.Process.Kill()
	}
	return nil
}

