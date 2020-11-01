// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package process

import (
	"os/exec"
	"syscall"
)

func Finalize(cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
	}
	return cmd.Wait()
}
