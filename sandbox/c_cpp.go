package sandbox

import (
	libseccomp "github.com/seccomp/libseccomp-golang"
	"os"
	"syscall"
)

func _cCppSeccompRules(_config *Config, allow_wirte_file bool) int {
	syscallsWhitelist := []int{
		syscall.SYS_READ, syscall.SYS_FSTAT,
		syscall.SYS_MMAP, syscall.SYS_MPROTECT,
		syscall.SYS_MUNMAP, syscall.SYS_UNAME,
		syscall.SYS_ARCH_PRCTL, syscall.SYS_BRK,
		syscall.SYS_ACCESS, syscall.SYS_EXIT_GROUP,
		syscall.SYS_CLOSE, syscall.SYS_READLINK,
		syscall.SYS_SYSINFO, syscall.SYS_WRITE,
		syscall.SYS_WRITEV, syscall.SYS_LSEEK,
		syscall.SYS_CLOCK_GETTIME, syscall.SYS_PREAD64,
	}
	filter, err := libseccomp.NewFilter(libseccomp.ActKillThread)
	if err != nil {
		return LOAD_SECCOMP_FAILED
	}
	for _, b := range syscallsWhitelist {
		if err := filter.AddRule(libseccomp.ScmpSyscall(b), libseccomp.ActAllow); err != nil {
			return LOAD_SECCOMP_FAILED
		}
	}
	defer filter.Release()
	// add extra rule for execve
	filter.AddRule(libseccomp.ScmpSyscall(syscall.SYS_EXECVE), libseccomp.ActAllow)
	if !allow_wirte_file {
		// do not allow "w" and "rw"
		// do not allow "w" and "rw" using open
		if err := filter.AddRuleConditional(
			syscall.SYS_OPEN,
			libseccomp.ActAllow,
			[]libseccomp.ScmpCondition{
				libseccomp.ScmpCondition{
					Argument: 1,
					Op:       libseccomp.CompareMaskedEqual,
					Operand1: uint64(os.O_WRONLY | os.O_RDWR),
					Operand2: uint64(os.O_WRONLY),
				},
			},
		); err != nil {
			return LOAD_SECCOMP_FAILED
		}
		// do not allow "w" and "rw" using openat
		if err := filter.AddRuleConditional(
			syscall.SYS_OPENAT,
			libseccomp.ActAllow,
			[]libseccomp.ScmpCondition{
				libseccomp.ScmpCondition{
					Argument: 2,
					Op:       libseccomp.CompareMaskedEqual,
					Operand1: uint64(os.O_WRONLY | os.O_RDWR),
					Operand2: 0,
				},
			},
		); err != nil {
			return LOAD_SECCOMP_FAILED
		}
	} else {
		if err := filter.AddRule(libseccomp.ScmpSyscall(syscall.SYS_OPEN), libseccomp.ActAllow); err != nil {
			return LOAD_SECCOMP_FAILED
		}

		if err := filter.AddRule(libseccomp.ScmpSyscall(syscall.SYS_DUP), libseccomp.ActAllow); err != nil {
			return LOAD_SECCOMP_FAILED
		}

		if err := filter.AddRule(libseccomp.ScmpSyscall(syscall.SYS_DUP2), libseccomp.ActAllow); err != nil {
			return LOAD_SECCOMP_FAILED
		}

		if err := filter.AddRule(libseccomp.ScmpSyscall(syscall.SYS_DUP3), libseccomp.ActAllow); err != nil {
			return LOAD_SECCOMP_FAILED
		}

	}

	if err := filter.Load(); err != nil {
		return LOAD_SECCOMP_FAILED
	}
	return 0
}

func cCppSeccompRules(_config *Config) int {
	return _cCppSeccompRules(_config, false)
}

func cCppFileIoSeccompRules(_config *Config) int {
	return _cCppSeccompRules(_config, true)
}
