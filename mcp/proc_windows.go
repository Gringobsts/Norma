//go:build windows

package mcp

import (
	"os"
	"os/exec"
	"strconv"
)

// setProcGroup is a no-op: taskkill /T walks the tree without a group.
func setProcGroup(*exec.Cmd) {}

// treeKill terminates the process and its descendants via taskkill.
func treeKill(p *os.Process) {
	if err := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(p.Pid)).Run(); err != nil {
		_ = p.Kill()
	}
}
