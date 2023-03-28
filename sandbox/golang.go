package sandbox

import (
	libseccomp "github.com/seccomp/libseccomp-golang"
	"os"
	"syscall"
)

func golangSeccompRules(_config *Config) int {
	syscallBlacklist := []int{
		syscall.SYS_SOCKET,
		syscall.SYS_FORK, syscall.SYS_VFORK,
		syscall.SYS_KILL,
		//syscall.SYS_EXECVEAT,
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
	defer filter.Release()
	// do not allow "w" and "rw" using open
	if err := filter.AddRuleConditional(
		syscall.SYS_OPEN,
		libseccomp.ActKillThread,
		[]libseccomp.ScmpCondition{
			libseccomp.ScmpCondition{
				Argument: 1,
				Op:       libseccomp.CompareMaskedEqual,
				Operand1: uint64(os.O_WRONLY),
				Operand2: uint64(os.O_WRONLY),
			}, {
				Argument: 1,
				Op:       libseccomp.CompareMaskedEqual,
				Operand1: uint64(os.O_RDWR),
				Operand2: uint64(os.O_RDWR),
			},
		},
	); err != nil {
		return LOAD_SECCOMP_FAILED
	}
	if err := filter.AddRuleConditional(
		syscall.SYS_OPENAT,
		libseccomp.ActKillThread,
		[]libseccomp.ScmpCondition{
			libseccomp.ScmpCondition{
				Argument: 2,
				Op:       libseccomp.CompareMaskedEqual,
				Operand1: uint64(os.O_WRONLY),
				Operand2: uint64(os.O_WRONLY),
			}, {
				Argument: 2,
				Op:       libseccomp.CompareMaskedEqual,
				Operand1: uint64(os.O_RDWR),
				Operand2: uint64(os.O_RDWR),
			},
		},
	); err != nil {
		return LOAD_SECCOMP_FAILED
	}
	if err := filter.Load(); err != nil {
		return LOAD_SECCOMP_FAILED
	}
	return 0
}
