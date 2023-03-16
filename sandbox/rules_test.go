package sandbox

import (
	"fmt"
	libseccomp "github.com/seccomp/libseccomp-golang"
	"os/exec"
	"reflect"
	"syscall"
	"testing"
	"unsafe"
)

func TestExecveRules(t *testing.T) {
	filter, err := libseccomp.NewFilter(libseccomp.ActAllow)
	if err != nil {
		return
	}
	a := "/bin/ls"
	_ = filter.AddRuleConditional(
		syscall.SYS_EXECVE,
		libseccomp.ActKillThread,
		[]libseccomp.ScmpCondition{
			libseccomp.ScmpCondition{
				Argument: 0,
				Op:       libseccomp.CompareNotEqual,
				Operand1: uint64((*reflect.StringHeader)(unsafe.Pointer(&a)).Data),
				Operand2: 0,
			},
		},
	)
	b := "/bin/cat"
	fmt.Println(uint64((*reflect.StringHeader)(unsafe.Pointer(&a)).Data))
	fmt.Println(uint64((*reflect.StringHeader)(unsafe.Pointer(&b)).Data))
	command := exec.Command("/bin/cat")
	command.Start()

}
