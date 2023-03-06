package sandbox

const (
	ArgsMaxNumber = 256
	EnvMaxNumber  = 256

	SUCCESS           = 0
	InvalidConfig     = -1
	ForkFailed        = -2
	PthreadFailed     = -3
	WaitFailed        = -4
	RootRequired      = -5
	LoadSeccompFailed = -6
	SetrlimitFailed   = -7
	Dup2Failed        = -8
	SetuidFailed      = -9
	ExecveFailed      = -10
	SpjError          = -11
)

type config struct {
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
	Args                 [ArgsMaxNumber]string
	Env                  [EnvMaxNumber]string
	logPath              string
	SeccompRuleName      string
	Uid                  int64
	Gid                  int64
}

const (
	WrongAnswer           = -1
	CpuTimeLimitExceeded  = 1
	RealTimeLimitExceeded = 2
	MemoryLimitExceeded   = 3
	RuntimeError          = 4
	SystemError           = 5
)

type result struct {
	CpuTime  int
	RealTime int
	Memory   int64
	Signal   int
	ExitCode int
	Error    int
	Result   int
}

func run(_config *config, _result *result) {

}
