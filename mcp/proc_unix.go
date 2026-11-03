//go:build !windows

package mcp

import (
	"os"
	"os/exec"
	"syscall"
)

// setProcGroup puts the server in its own process group so Close can kill the
// whole tree (e.g. npx → node).
func setProcGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true
}

// treeKill terminates the process group, falling back to the lone process.
func treeKill(p *os.Process) {
	if err := syscall.Kill(-p.Pid, syscall.SIGKILL); err != nil {
		_ = p.Kill()
	}
}
