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
	flag.IntVar(&max_cpu_time, "max_cpu_time", 0, "Max CPU Time (ms)")
	flag.IntVar(&max_real_time, "max_real_time", 0, "Max Real Time (ms)")
	flag.IntVar(&max_memory, "max_memory", 0, "Max Memory (byte)")
	flag.IntVar(&max_stack, "max_stack", 0, "Max Stack (byte, default 16M)")
	flag.IntVar(&max_process_number, "max_process_number", 0, "Max Process Number")
	flag.IntVar(&max_output_size, "max_output_size", 0, "Max Output Size (byte)")
	flag.IntVar(&memory_limit_check_only, "memory_limit_check_only", 0, "only check memory usage, do not setrlimit (default False)")

	flag.StringVar(&exe_path, "exe_path", "", "Exe Path")
	flag.StringVar(&input_file, "input_file", "", "Input Path")
	flag.StringVar(&output_file, "output_file", "", "Output Path")
	flag.StringVar(&error_file, "error_file", "", "Error Path")
	flag.StringVar(&log_path, "log_path", "", "Log Path")
	flag.StringVar(&seccomp_rule_name, "seccomp_rule_name", "", "Seccomp Rule Name")

	flag.IntVar(&uid, "uid", 0, "UID (default 65534)")
	flag.IntVar(&gid, "gid", 0, "GID (default 65534)")
	flag.Var(&args, "args", "Arg")
	flag.Var(&env, "env", "Env")
}

func main() {
	flag.Parse()
	fmt.Println(args)
	fmt.Println(env)
	sandbox.LogOpen("")
}
