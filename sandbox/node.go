package sandbox

import (
	libseccomp "github.com/seccomp/libseccomp-golang"
	"syscall"
)

const (
	SYS_EXECVEAT = 281
)

func nodeSeccompRules(_config *Config) int {
	syscallBlacklist := []int{
		syscall.SYS_SOCKET,
		syscall.SYS_FORK, syscall.SYS_VFORK,
		syscall.SYS_KILL,
		SYS_EXECVEAT,
	}
	filter, err := libseccomp.NewFilter(libseccomp.ActAllow)
	if err != nil {
		return LOAD_SECCOMP_FAILED
	}
	for _, b := range syscallBlacklist {
		if err := filter.AddRule(libseccomp.ScmpSyscall(b), libseccomp.ActKillThread); err != nil {
			return LOAD_SECCOMP_FAILED
		}
	}
	if err := filter.Load(); err != nil {
		return LOAD_SECCOMP_FAILED
	}
	return 0
}
