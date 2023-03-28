package sandbox

/*
 int golang_seccomp_rules(const char*);
*/
import "C"

func golang_seccomp_rules(_config *Config) int {
	return int(C.golang_seccomp_rules(C.CString(_config.ExePath)))
}
