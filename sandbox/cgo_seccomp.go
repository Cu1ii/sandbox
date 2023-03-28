package sandbox

/*
//#include <seccomp.h>
//
//int seccomp_rule_add_extra_execve(const char* exe_path) {
//	const int LOAD_SECCOMP_FAILED = -6;
//	const int SUCCESS = 0;
//#define SCMP_SYS(execve) 59
//	// add extra rule for execve
//	if (seccomp_rule_add(ctx, SCMP_ACT_ALLOW, SCMP_SYS(execve), 1, SCMP_A0(SCMP_CMP_EQ, (scmp_datum_t)(exe_path))) != 0) {
//	        return LOAD_SECCOMP_FAILED;
//    }
//	return SUCCESS;
//}

*/
import "C"
import (
	seccomp "github.com/seccomp/libseccomp-golang"
)

func addSysExecveRule(_config *Config, filter *seccomp.ScmpFilter) {
	//res := C.seccomp_rule_add_extra_execve(C.CString(_config.ExePath))
	//fmt.Printf("%T %v", res, res)
}
