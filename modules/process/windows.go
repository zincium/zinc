// +build windows

package process

import "os/exec"

// Finalize todo
func Finalize(cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}
	if cmd.Process != nil && cmd.Process.Pid > 0 {
		_ = cmd.Process.Kill()
	}
	return cmd.Wait()
}
