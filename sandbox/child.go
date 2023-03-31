package sandbox

import (
	"fmt"
	"os"
	"syscall"
	// "strconv"
)

var inputFile = (*os.File)(nil)
var outputFile = (*os.File)(nil)
var errorFile = (*os.File)(nil)

func childErrorExit(logfile *os.File, code int) {
	msg := fmt.Sprintf("Error: System errno: %d ", code)
	LogFatal(logfile, msg)
	if outputFile == errorFile {
		CloseFile(outputFile)
	} else {
		CloseFile(outputFile)
		CloseFile(inputFile)
	}
	syscall.Kill(syscall.Getpid(), syscall.SIGUSR1)
	os.Exit(EXIT_FAILURE)
}

func childProcess(logfile *os.File, _config *Config) {

	var err error

	if _config.MaxStack != UNLIMITED {
		maxStack := syscall.Rlimit{
			Cur: uint64(_config.MaxStack),
			Max: uint64(_config.MaxStack),
		}
		err := syscall.Setrlimit(syscall.RLIMIT_STACK, &maxStack)
		if err != nil {
			childErrorExit(logfile, SETRLIMIT_FAILED)
			//os.Exit(SETRLIMIT_FAILED)
		}
	}
	// set memory limit
	// if memory_limit_check_only == 0, we only check memory usage number,
	// because setrlimit(maxrss) will cause some crash issues
	if _config.MemoryLimitCheckOnly == 0 {
		maxMemory := syscall.Rlimit{
			Cur: uint64(_config.MaxMemory * 2),
			Max: uint64(_config.MaxMemory * 2),
		}
		err := syscall.Setrlimit(syscall.RLIMIT_AS, &maxMemory)
		if err != nil {
			childErrorExit(logfile, SETRLIMIT_FAILED)
			//os.Exit(SETRLIMIT_FAILED)
		}
	}

	// set cpu time limit (in seconds)
	if _config.MaxCpuTime == UNLIMITED {
		maxCpuTime := syscall.Rlimit{
			Cur: uint64((_config.MaxCpuTime + 1000) / 1000),
			Max: uint64((_config.MaxCpuTime + 1000) / 1000),
		}
		err := syscall.Setrlimit(syscall.RLIMIT_CPU, &maxCpuTime)
		if err != nil {
			childErrorExit(logfile, SETRLIMIT_FAILED)
			//                      os.Exit(SETRLIMIT_FAILED)
		}
	}

	// https://go.dev/misc/cgo/test/issue18146.go
	const SYS_RLIMIT_NPROC = 6

	// set process number limit
	if _config.MaxProcessNumber != UNLIMITED {
		maxProcessNumber := syscall.Rlimit{
			Cur: uint64(_config.MaxProcessNumber),
			Max: uint64(_config.MaxProcessNumber),
		}
		err := syscall.Setrlimit(SYS_RLIMIT_NPROC, &maxProcessNumber)
		if err != nil {
			childErrorExit(logfile, SETRLIMIT_FAILED)
			//                      os.Exit(SETRLIMIT_FAILED)
		}
	}
	// set output size limit
	if _config.MaxOutputSize != UNLIMITED {
		maxOutputSize := syscall.Rlimit{
			Cur: uint64(_config.MaxOutputSize),
			Max: uint64(_config.MaxOutputSize),
		}
		err := syscall.Setrlimit(syscall.RLIMIT_FSIZE, &maxOutputSize)
		if err != nil {
			childErrorExit(logfile, SETRLIMIT_FAILED)
			// os.Exit(SETRLIMIT_FAILED)
		}
	}

	if _config.InputPath != "" {
		inputFile, err = os.Open(_config.InputPath)
		if err != nil {
			childErrorExit(logfile, DUP2_FAILED)
			// os.Exit(DUP2_FAILED)
		}
		// redirect file -> stdin
		// On success, these system calls return the new descriptor.
		// On error, -1 is returned, and errno is set appropriately.
		fmt.Println("child input_file: ", inputFile.Name())
		if err := syscall.Dup2(int(inputFile.Fd()), int(os.Stdin.Fd())); err != nil {
			childErrorExit(logfile, DUP2_FAILED)
			// os.Exit(DUP2_FAILED)
		}
	}

	if _config.OutputPath != "" {
		outputFile, err = os.OpenFile(_config.OutputPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			childErrorExit(logfile, DUP2_FAILED)
			//os.Exit(DUP2_FAILED)
		}
		fmt.Println("child output_path: ", outputFile.Name())
		if err := syscall.Dup2(int(outputFile.Fd()), int(os.Stdout.Fd())); err != nil {
			childErrorExit(logfile, DUP2_FAILED)
			// os.Exit(DUP2_FAILED)
		}
	}
	if _config.ErrorPath != "" {
		// if outfile and error_file is the same path, we use the same file pointer
		if _config.OutputPath != "" && _config.ErrorPath == _config.OutputPath {
			errorFile = outputFile
		} else {

			errorFile, err = os.OpenFile(_config.ErrorPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				childErrorExit(logfile, DUP2_FAILED)
				// os.Exit(DUP2_FAILED)
			}
		}
		// redirect stderr -> file
		if err := syscall.Dup2(int(errorFile.Fd()), int(os.Stderr.Fd())); err != nil {
			childErrorExit(logfile, DUP2_FAILED)
			// os.Exit(DUP2_FAILED)
		}
	}

	if _config.Gid != -1 {
		if err := syscall.Setegid(int(_config.Gid)); err != nil {
			childErrorExit(logfile, SETUID_FAILED)
			// os.Exit(SETUID_FAILED)
		}
	}

	if _config.Uid != -1 {
		if err := syscall.Seteuid(int(_config.Uid)); err != nil {
			childErrorExit(logfile, SETUID_FAILED)
			// os.Exit(SETUID_FAILED)
		}
	}

	if _config.SeccompRuleName != "" {
		if "c_cpp" == _config.SeccompRuleName {
			if c_cpp_file_io_seccomp_rules(_config) != SUCCESS {
				childErrorExit(logfile, LOAD_SECCOMP_FAILED)
				// os.Exit(LOAD_SECCOMP_FAILED)
			}
		} else if "c_cpp_file_io" == _config.SeccompRuleName {
			if c_cpp_file_io_seccomp_rules(_config) != SUCCESS {
				childErrorExit(logfile, LOAD_SECCOMP_FAILED)
				// os.Exit(LOAD_SECCOMP_FAILED)
			}
		} else if "general" == _config.SeccompRuleName {
			if general_seccomp_rules(_config) != SUCCESS {
				childErrorExit(logfile, LOAD_SECCOMP_FAILED)
				// os.Exit(LOAD_SECCOMP_FAILED)
			}
		} else if "golang" == _config.SeccompRuleName {
			if golang_seccomp_rules(_config) != SUCCESS {
				childErrorExit(logfile, LOAD_SECCOMP_FAILED)
				// os.Exit(LOAD_SECCOMP_FAILED)
			}
		} else if "node" == _config.SeccompRuleName {
			if node_seccomp_rules(_config) != SUCCESS {
				childErrorExit(logfile, LOAD_SECCOMP_FAILED)
				// os.Exit(LOAD_SECCOMP_FAILED)
			}
		} else {
			childErrorExit(logfile, LOAD_SECCOMP_FAILED)
			// os.Exit(LOAD_SECCOMP_FAILED)
		}
	}
	if err := syscall.Exec(_config.ExePath, _config.Args, _config.Env); err != nil {
		childErrorExit(logfile, EXECVE_FAILED)
		// os.Exit(EXECVE_FAILED)
	}
	childErrorExit(logfile, EXECVE_FAILED)
	//os.StartProcess(_config.ExePath)
}
