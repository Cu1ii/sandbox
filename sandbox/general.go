package sandbox

import (
	libseccomp "github.com/seccomp/libseccomp-golang"
	"os"
	"syscall"
)

func generalSeccompRules(_config *Config) int {
	syscallBlacklist := []int{
		syscall.SYS_CLONE,
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
	// use libseccomp.ActKillThread for socket, python will be killed immediately
	if err := filter.AddRule(syscall.SYS_SOCKET, libseccomp.ActErrno.SetReturnCode(int16(syscall.EACCES))); err != nil {
		return LOAD_SECCOMP_FAILED
	}

	// add extra rule for execve
	//if err := filter.AddRuleConditional(
	//	syscall.SYS_EXECVE,
	//	libseccomp.ActKillThread,
	//	[]libseccomp.ScmpCondition{
	//		libseccomp.ScmpCondition{
	//			Argument: 0,
	//			Op:       libseccomp.CompareNotEqual,
	//			Operand1: uint64((*reflect.StringHeader)(unsafe.Pointer(&_config.ExePath)).Data),
	//			Operand2: 0,
	//		},
	//	},
	//); err != nil {
	//	return LOAD_SECCOMP_FAILED
	//}

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
	// do not allow "w" and "rw" using openat
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
