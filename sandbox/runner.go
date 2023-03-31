package sandbox

import (
	"errors"
	"os"
	"strconv"
	"syscall"
	"time"
)

const (
	UNLIMITED = -1

	ArgsMaxNumber = 256
	EnvMaxNumber  = 256

	SUCCESS             = 0
	INVALID_CONFIG      = -1
	ForkFailed          = -2
	PTHREAD_FAILED      = -3
	WaitFailed          = -4
	ROOT_REQUIRED       = -5
	LOAD_SECCOMP_FAILED = -6
	SETRLIMIT_FAILED    = -7
	DUP2_FAILED         = -8
	SETUID_FAILED       = -9
	EXECVE_FAILED       = -10
	SPJ_ERROR           = -11

	EXIT_FAILURE = 1
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
	WRONG_ANSWER             = -1
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

var logfile = (*os.File)(nil)

func Run(_config *Config, _result *Result) error {
	logfile = LogOpen(_config.LogPath)
	defer CloseFile(logfile)
	initResult(_result)
	uid := os.Geteuid()
	if uid != 0 {
		LogDebug(logfile, strconv.Itoa(ROOT_REQUIRED))
		_result.Result = ROOT_REQUIRED
		return errors.New("ROOT_REQUIRED")
		// os.Exit(ROOT_REQUIRED)
	}
	if (_config.MaxCpuTime < 1 && _config.MaxCpuTime != UNLIMITED) ||
		(_config.MaxRealTime < 1 && _config.MaxRealTime != UNLIMITED) ||
		(_config.MaxStack < 1) ||
		(_config.MaxMemory < 1 && _config.MaxMemory != UNLIMITED) ||
		(_config.MaxProcessNumber < 1 && _config.MaxProcessNumber != UNLIMITED) ||
		(_config.MaxOutputSize < 1 && _config.MaxOutputSize != UNLIMITED) {
		LogDebug(logfile, strconv.Itoa(INVALID_CONFIG))
		_result.Result = INVALID_CONFIG
		return errors.New("INVALID_CONFIG")
		// os.Exit(INVALID_CONFIG)
	}

	start := time.Now()
	childPid, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		LogDebug(logfile, strconv.Itoa(ForkFailed))
		_result.Result = ForkFailed
		return errors.New("ForkFailed")
		// os.Exit(ForkFailed)
	}
	process, err := os.FindProcess(int(childPid))
	if err != nil {
		LogDebug(logfile, strconv.Itoa(ForkFailed))
		_result.Result = ForkFailed
		return errors.New("ForkFailed")
		//  LogDebug(logfile, "get child process failed")
		// os.Exit(ForkFailed)
	}
	if int(childPid) < 0 {
		LogDebug(logfile, strconv.Itoa(ForkFailed))
		_result.Result = ForkFailed
		return errors.New("ForkFailed")
		// os.Exit(ForkFailed)
	} else if int(childPid) == 0 {
		//LogDebug(logfile, "child process")
		//time.Sleep(12)
		childProcess(logfile, _config)
	} else if int(childPid) > 0 {
		if _config.MaxRealTime != UNLIMITED {
			timeout := _config.MaxRealTime
			go func() {
				_ = timeKiller(process.Pid, timeout)
				return
			}()
		}

		rusage := syscall.Rusage{}
		var status syscall.WaitStatus
		_, err := syscall.Wait4(int(childPid), &status, syscall.WSTOPPED, &rusage)

		if err != nil {
			LogDebug(logfile, strconv.Itoa(WaitFailed))
			_ = killProcess(process.Pid)
			// os.Exit(WaitFailed)
		}
		end := time.Now()
		_result.RealTime = int(end.Unix() - start.Unix())
		if _config.MaxRealTime != UNLIMITED {

		}
		//fmt.Println(status.Signal())
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
