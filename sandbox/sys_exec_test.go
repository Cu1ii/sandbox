package sandbox

import (
	"fmt"
	"github.com/seccomp/libseccomp-golang"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestSysExec(t *testing.T) {
	cmd := "/bin/ls"
	args := []string{"ls", "-l", "/"}
	env := os.Environ()

	err := syscall.Exec(cmd, args, env)
	if err != nil {
		panic(err)
	}

	fmt.Println("This line will never be executed")
}

func TestSec(t *testing.T) {
	// 创建seccomp过滤器
	filter, err := seccomp.NewFilter(seccomp.ActKillThread)
	if err != nil {
		panic(err)
	}
	defer filter.Release()

	// 添加规则，允许读取当前进程ID
	filter.AddRule(syscall.SYS_GETPID, seccomp.ActAllow)

	// 添加规则，限制execve系统调用只能执行指定的可执行文件路径

	filter.AddRuleConditional(syscall.SYS_EXECVE, seccomp.ActKillThread, []seccomp.ScmpCondition{
		seccomp.ScmpCondition{
			Argument: 0,
			Op:       seccomp.CompareNotEqual,
		},
	})
	// 加载过滤器
	if err := filter.Load(); err != nil {
		panic(err)
	}

	// 执行ls命令
	cmd := &exec.Cmd{
		Path: "/bin/ls",
		Args: []string{"/bin/ls", "-l"},
	}
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

func TestCgoCcpp(t *testing.T) {
	res := c_cpp_seccomp_rules(&Config{})
	fmt.Printf("%T %v", res, res)
}

func TestCgoGeneral(t *testing.T) {
	res := general_seccomp_rules(&Config{})
	fmt.Printf("%T %v", res, res)
}

func TestCgoGolang(t *testing.T) {
	res := golang_seccomp_rules(&Config{})
	fmt.Printf("%T %v", res, res)
}

func TestCgoNode(t *testing.T) {
	res := node_seccomp_rules(&Config{})
	fmt.Printf("%T %v", res, res)
}

func TestSayHello(t *testing.T) {
	sayHello()
}

func TestSysE(t *testing.T) {
	if err := syscall.Exec("ls", []string{"-a", "/home/cu1"}, []string{}); err != nil {
		//childErrorExit(logfile, EXECVE_FAILED)
		fmt.Println("FAILED")
		// os.Exit(EXECVE_FAILED)
	}
}
