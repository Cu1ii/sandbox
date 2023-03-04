package main

import (
	"flag"
	"fmt"
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
	flag.IntVar(&max_cpu_time, "max_cpu_time", 0, "(ms): max cpu time this process can cost, -1 for unlimited")
	flag.IntVar(&max_real_time, "max_real_time", 0, "(ms): max time this process can run, -1 for unlimited")
	flag.IntVar(&max_memory, "max_memory", 0, "(byte): max size of the process' virtual memory (address space), -1 for unlimited")
	flag.IntVar(&max_stack, "max_stack", 0, "(byte): max size of the process' stack size")
	flag.IntVar(&max_process_number, "max_process_number", 0, "max number of processes that can be created for the real user id of the calling process, -1 for unlimited")
	flag.IntVar(&max_output_size, "max_output_size", 0, "(byte): max size of data this process can output to stdout, stderr and file, -1 for unlimited")
	flag.IntVar(&memory_limit_check_only, "memory_limit_check_only", 0, "if this value equals 0, we will only check memory usage number, because setrlimit(maxrss) will cause some crash issues")
	flag.StringVar(&exe_path, "exe_path", "", "path of file to run")
	flag.StringVar(&input_file, "input_file", "", "redirect content of this file to process's stdin")
	flag.StringVar(&output_file, "output_file", "", "redirect process's stdout to this file")
	flag.StringVar(&error_file, "error_file", "", "redirect process's stderr to this file")
	flag.StringVar(&log_path, "log_path", "", "judger log path")
	flag.StringVar(&seccomp_rule_name, "seccomp_rule_name", "", "(string or NULL): seccomp rules used to limit process system calls. Name is used to call corresponding functions.")
	flag.IntVar(&uid, "uid", 0, "user to run this process")
	flag.IntVar(&gid, "gid", 0, "user group this process belongs to")
	flag.Var(&args, "args", "(string array terminated by NULL): arguments to run this process")
	flag.Var(&env, "env", "(string array terminated by NULL): environment variables this process can get")
}

func main() {
	flag.Parse()
	fmt.Println(args)
}
