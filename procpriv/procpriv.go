// Package procpriv optionally runs tool subprocesses as an unprivileged user.
//
// Set NORMA_TOOL_UID (and NORMA_TOOL_GID) to enable it. The host process keeps
// its own privileges, so files it owns (secrets, its /proc environ) stay out of
// reach of the commands an agent runs. Unset: everything behaves as before.
//
//	NORMA_TOOL_UID / NORMA_TOOL_GID  numeric ids for tool processes
//	NORMA_TOOL_HOME                  HOME for tool processes (default /tmp)
//	NORMA_TOOL_WORKROOT              working dirs under this root are chowned to
//	                                 the tool user so it can write there
package procpriv

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type ids struct {
	uid, gid       int
	home, workRoot string
}

var load = sync.OnceValue(func() *ids {
	uid, err := strconv.Atoi(strings.TrimSpace(os.Getenv("NORMA_TOOL_UID")))
	if err != nil || uid <= 0 {
		return nil
	}
	gid, err := strconv.Atoi(strings.TrimSpace(os.Getenv("NORMA_TOOL_GID")))
	if err != nil || gid <= 0 {
		gid = uid
	}
	home := os.Getenv("NORMA_TOOL_HOME")
	if home == "" {
		home = "/tmp"
	}
	root := os.Getenv("NORMA_TOOL_WORKROOT")
	if root != "" {
		root = filepath.Clean(root)
	}
	return &ids{uid: uid, gid: gid, home: home, workRoot: root}
})

// Enabled reports whether tool processes are dropped to another user.
func Enabled() bool { return load() != nil }

// Chown hands path to the tool user when enabled and path is under the work
// root, so files the host writes for the agent stay editable by its commands.
func Chown(path string) {
	p := load()
	if p == nil || p.workRoot == "" || !within(p.workRoot, path) {
		return
	}
	fi, err := os.Lstat(path)
	if err != nil || fi.Mode()&os.ModeSymlink != 0 {
		return
	}
	if os.Lchown(path, p.uid, p.gid) == nil {
		return
	}
	// Some bind mounts (Docker Desktop file sharing) refuse chown: open the
	// mode instead so the tool user can still write its own work dir.
	mode := os.FileMode(0o666)
	if fi.IsDir() {
		mode = 0o777
	}
	_ = os.Chmod(path, mode)
}

func within(root, path string) bool {
	rel, err := filepath.Rel(root, filepath.Clean(path))
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// withHome replaces HOME in env (nil = inherit os.Environ()).
func withHome(env []string, home string) []string {
	if env == nil {
		env = os.Environ()
	}
	out := make([]string, 0, len(env)+1)
	for _, kv := range env {
		if !strings.HasPrefix(kv, "HOME=") {
			out = append(out, kv)
		}
	}
	return append(out, "HOME="+home)
}
