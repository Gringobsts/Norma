package procpriv

import "os/exec"

// Apply is a no-op on Windows.
func Apply(*exec.Cmd) {}
