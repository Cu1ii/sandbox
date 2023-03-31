package main

import (
	"flag"
	"fmt"
	"sandbox/sandbox"
)

type Args []string
type Env []string

func (a *Args) String() string {
	return fmt.Sprint(*a)
}

func (a *Args) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func (e *Env) String() string {
	return fmt.Sprint(*e)
}

func (e *Env) Set(value string) error {
	*e = append(*e, value)
	return nil
}

var (
	max_cpu_time            int    // (ms): max cpu time this process can cost, -1 for unlimited
	max_real_time           int    // (ms): max time this process can run, -1 for unlimited
	max_memory              int    // (byte): max size of the process' virtual memory (address space), -1 for unlimited
	max_stack               int    // (byte): max size of the process' stack size
	max_process_number      int    // max number of processes that can be created for the real user id of the calling process, -1 for unlimited
	max_output_size         int    // (byte): max size of data this process can output to stdout, stderr and file, -1 for unlimited
	memory_limit_check_only int    // if this value equals 0, we will only check memory usage number, because setrlimit(maxrss) will cause some crash issues
	exe_path                string // path of file to run
	input_file              string // redirect content of this file to process's stdin
	output_file             string // redirect process's stdout to this file
	error_file              string // redirect process's stderr to this file
	args                    Args   // (string array terminated by NULL): arguments to run this process
	env                     Env    // (string array terminated by NULL): environment variables this process can get
	log_path                string // judger log path
	seccomp_rule_name       string // (string or NULL): seccomp rules used to limit process system calls. Name is used to call corresponding functions.
	uid                     int    // user to run this process
	gid                     int    // user group this process belongs to
)

func init() {
	flag.IntVar(&max_cpu_time, "max_cpu_time", -1, "Max CPU Time (ms)")
	flag.IntVar(&max_real_time, "max_real_time", -1, "Max Real Time (ms)")
	flag.IntVar(&max_memory, "max_memory", -1, "Max Memory (byte)")
	flag.IntVar(&max_stack, "max_stack", -1, "Max Stack (byte, default 16M)")
	flag.IntVar(&max_process_number, "max_process_number", -1, "Max Process Number")
	flag.IntVar(&max_output_size, "max_output_size", -1, "Max Output Size (byte)")
	flag.IntVar(&memory_limit_check_only, "memory_limit_check_only", -1, "only check memory usage, do not setrlimit (default False)")

	flag.StringVar(&exe_path, "exe_path", "", "Exe Path")
	flag.StringVar(&input_file, "input_file", "", "Input Path")
	flag.StringVar(&output_file, "output_file", "", "Output Path")
	flag.StringVar(&error_file, "error_path", "", "Error Path")
	flag.StringVar(&log_path, "log_path", "", "Log Path")
	flag.StringVar(&seccomp_rule_name, "seccomp_rule_name", "", "Seccomp Rule Name")

	flag.IntVar(&uid, "uid", 65534, "UID (default 65534)")
	flag.IntVar(&gid, "gid", 65534, "GID (default 65534)")
	flag.Var(&args, "args", "Arg")
	flag.Var(&env, "env", "Env")
}

func main() {
	flag.Parse()
	_config := sandbox.Config{}
	if max_cpu_time > -1 {
		_config.MaxCpuTime = max_cpu_time
	} else {
		_config.MaxCpuTime = sandbox.UNLIMITED
	}

	if max_real_time > -1 {
		_config.MaxRealTime = max_real_time
	} else {
		_config.MaxRealTime = sandbox.UNLIMITED
	}

	if max_memory > -1 {
		_config.MaxMemory = int64(max_memory)
	} else {
		_config.MaxMemory = sandbox.UNLIMITED
	}

	if memory_limit_check_only > -1 {
		_config.MemoryLimitCheckOnly = memory_limit_check_only
	} else {
		_config.MemoryLimitCheckOnly = 0
	}

	if max_stack > -1 {
		_config.MaxStack = int64(max_stack)
	} else {
		_config.MaxStack = 16 * 1024 * 1024
	}

	if max_process_number > -1 {
		_config.MaxProcessNumber = max_process_number
	} else {
		_config.MaxProcessNumber = sandbox.UNLIMITED
	}

	if max_output_size > -1 {
		_config.MaxOutputSize = int64(max_output_size)
	} else {
		_config.MaxOutputSize = sandbox.UNLIMITED
	}

	_config.ExePath = exe_path

	if input_file != "" {
		_config.InputPath = input_file
	} else {
		_config.InputPath = "/dev/stdin"
	}

	if output_file != "" {
		_config.OutputPath = output_file
	} else {
		_config.OutputPath = "/dev/stdout"
	}

	if error_file != "" {
		_config.ErrorPath = error_file
	} else {
		_config.ErrorPath = "/dev/stderr"
	}

	if len(args) > 0 {
		copy(_config.Args, args)
	}

	if len(env) > 0 {
		copy(_config.Env, env)
	}

	if log_path != "" {
		_config.LogPath = log_path
	} else {
		_config.LogPath = "judger.log"
	}

	if seccomp_rule_name != "" {
		_config.SeccompRuleName = seccomp_rule_name
	}

	if uid > -1 {
		_config.Uid = int64(uid)
	} else {
		_config.Uid = 65534
	}

	if gid > -1 {
		_config.Gid = int64(gid)
	} else {
		_config.Gid = 65534
	}
	fmt.Printf("cpu_time = %d, real_time = %d, memory = %d, stack = %d, process = %d, output_size = %d",
		_config.MaxCpuTime, _config.MaxRealTime, _config.MaxMemory, _config.MaxStack, _config.MaxProcessNumber, _config.MaxOutputSize)
	_result := sandbox.Result{}
	sandbox.Run(&_config, &_result)
	str := "{\n" +
		"    \"cpu_time\": %d,\n" +
		"    \"real_time\": %d,\n" +
		"    \"memory\": %d,\n" +
		"    \"signal\": %d,\n" +
		"    \"exit_code\": %d,\n" +
		"    \"error\": %d,\n" +
		"    \"result\": %d\n" +
		"}"
	fmt.Printf(str,
		_result.CpuTime,
		_result.RealTime,
		_result.Memory,
		_result.Signal,
		_result.ExitCode,
		_result.Error,
		_result.Result)
}
