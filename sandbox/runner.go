package sandbox

import (
	"os"
	"syscall"
	"time"
)

const (
	UNLIMITED = -1

	ArgsMaxNumber = 256
	EnvMaxNumber  = 256

	SUCCESS             = 0
	InvalidConfig       = -1
	ForkFailed          = -2
	PthreadFailed       = -3
	WaitFailed          = -4
	RootRequired        = -5
	LOAD_SECCOMP_FAILED = -6
	SETRLIMIT_FAILED    = -7
	DUP2_FAILED         = -8
	SETUID_FAILED       = -9
	EXECVE_FAILED       = -10
	SpjError            = -11
)

type Config struct {
	MaxCpuTime           int
	MaxRealTime          int
	MaxMemory            int64
	MaxStack             int64
	MaxProcessNumber     int
	MaxOutputSize        int64
	MemoryLimitCheckOnly int
	ExePath              string
	InputPath            string
	OutputPath           string
	ErrorPath            string
	Args                 []string
	Env                  []string
	LogPath              string
	SeccompRuleName      string
	Uid                  int64
	Gid                  int64
}

const (
	WrongAnswer              = -1
	CPU_TIME_LIMIT_EXCEEDED  = 1
	REAL_TIME_LIMIT_EXCEEDED = 2
	MEMORY_LIMIT_EXCEEDED    = 3
	RUNTIME_ERROR            = 4
	SYSTEM_ERROR             = 5
)

type Result struct {
	CpuTime  int
	RealTime int
	Memory   int64
	Signal   int
	ExitCode int
	Error    int
	Result   int
}

func initResult(_result *Result) {
	_result.Result = SUCCESS
	_result.Error = SUCCESS
	_result.CpuTime = 0
	_result.RealTime = 0
	_result.Signal = 0
	_result.ExitCode = 0
	_result.Memory = 0
}

func Run(_config *Config, _result *Result) error {
	logfile := LogOpen(_config.LogPath)
	defer CloseFile(logfile)
	initResult(_result)
	uid := os.Geteuid()
	if uid != 0 {
		os.Exit(RootRequired)
	}
	if (_config.MaxCpuTime < 1 && _config.MaxCpuTime != UNLIMITED) ||
		(_config.MaxRealTime < 1 && _config.MaxRealTime != UNLIMITED) ||
		(_config.MaxStack < 1) ||
		(_config.MaxMemory < 1 && _config.MaxMemory != UNLIMITED) ||
		(_config.MaxProcessNumber < 1 && _config.MaxProcessNumber != UNLIMITED) ||
		(_config.MaxOutputSize < 1 && _config.MaxOutputSize != UNLIMITED) {
		os.Exit(InvalidConfig)
	}

	start := time.Now()
	childPid, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		os.Exit(ForkFailed)
	}
	process, err := os.FindProcess(int(childPid))
	if err != nil {
		LogDebug(logfile, "get child process failed")
		os.Exit(ForkFailed)
	}
	if childPid < 0 {
		os.Exit(ForkFailed)
	} else if childPid == 0 {
		childProcess(logfile, _config)
	} else if childPid > 0 {
		if _config.MaxRealTime != UNLIMITED {
			timeout := _config.MaxRealTime
			go func() {
				_ = timeKiller(process.Pid, timeout)
			}()
		}

		rusage := syscall.Rusage{}
		var status syscall.WaitStatus
		_, err := syscall.Wait4(process.Pid, &status, syscall.WSTOPPED, &rusage)
		if err != nil {
			_ = killProcess(process.Pid)
			os.Exit(WaitFailed)
		}
		end := time.Now()
		_result.RealTime = end.Second() - start.Second()
		if _config.MaxRealTime != UNLIMITED {

		}
		if status.Signal() != 0 {
			_result.Signal = int(status.Signal())
		}

		if _result.Signal == int(syscall.SIGUSR1) {
			_result.Result = SYSTEM_ERROR
		} else {
			_result.ExitCode = status.ExitStatus()
			_result.CpuTime = int(rusage.Utime.Sec*1000 + rusage.Utime.Usec/1000)
			_result.Memory = rusage.Maxrss * 1024
			if _result.ExitCode != 0 {
				_result.Result = RUNTIME_ERROR
			}
			if _result.Signal == int(syscall.SIGSEGV) {
				if _config.MaxMemory != UNLIMITED && _result.Memory > _config.MaxMemory {
					_result.Result = MEMORY_LIMIT_EXCEEDED
				} else {
					_result.Result = RUNTIME_ERROR
				}
			} else {
				if _result.Signal != 0 {
					_result.Result = RUNTIME_ERROR
				}
				if _config.MaxMemory != UNLIMITED && _result.Memory > _config.MaxMemory {
					_result.Result = MEMORY_LIMIT_EXCEEDED
				}
				if _config.MaxRealTime != UNLIMITED && _result.RealTime > _config.MaxRealTime {
					_result.Result = REAL_TIME_LIMIT_EXCEEDED
				}
				if _config.MaxCpuTime != UNLIMITED && _result.CpuTime > _config.MaxCpuTime {
					_result.Result = CPU_TIME_LIMIT_EXCEEDED
				}
			}
		}
	}

	return nil
}
