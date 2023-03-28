package sandbox

/*
 int node_seccomp_rules(const char* );
*/
import "C"

func node_seccomp_rules(_config *Config) int {
	return int(C.node_seccomp_rules(C.CString(_config.ExePath)))
}
