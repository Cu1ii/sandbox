package sandbox

/*
 int general_seccomp_rules(const char*);
*/
import "C"

func general_seccomp_rules(_config *Config) int {
	return int(C.general_seccomp_rules(C.CString(_config.ExePath)))
}
