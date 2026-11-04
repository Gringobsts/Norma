//go:build !windows

package procpriv

import (
	"os/exec"
	"syscall"
)

// Apply makes cmd run as the tool user when enabled; call it before Start.
// The working dir is chowned to that user when it sits under the work root.
func Apply(cmd *exec.Cmd) {
	p := load()
	if p == nil {
		return
	}
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(p.uid), Gid: uint32(p.gid)}
	cmd.Env = withHome(cmd.Env, p.home)
	if cmd.Dir != "" {
		Chown(cmd.Dir)
	}
}
