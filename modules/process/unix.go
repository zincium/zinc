//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package process

import (
	"os/exec"
	"syscall"
)

func Finalize(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil || cmd.Process.Pid <= 0 {
		return nil
	}
	if cmd.SysProcAttr != nil {
		if cmd.SysProcAttr.Setpgid {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
	}
	syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	return cmd.Wait()
}
