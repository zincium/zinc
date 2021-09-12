package process

import (
	"errors"
	"os/exec"
)

// TODO

func ExitCode(err error) (int32, bool) {
	if err == nil {
		return 0, false
	}
	var exitError *exec.ExitError
	if ok := errors.As(err, &exitError); !ok {
		return 0, false
	}
	return int32(exitError.ExitCode()), true
}
