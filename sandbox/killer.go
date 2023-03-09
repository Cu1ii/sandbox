package sandbox

import (
	"os"
	"time"
)

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

func timeKiller(pid int, timeout int) error {
	_ = <-time.After(time.Duration((timeout+1000)/1000) * time.Second)
	err := killProcess(pid)
	if err != nil {
		return err
	}
	return nil
}
