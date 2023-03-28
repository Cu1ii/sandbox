package sandbox

/*
int c_cpp_seccomp_rules(const char* , int);
*/
import "C"

func c_cpp_seccomp_rules(_config *Config) int {
	return int(C.c_cpp_seccomp_rules(C.CString(_config.ExePath), C.int(0)))
}

func c_cpp_file_io_seccomp_rules(_config *Config) int {
	return int(C.c_cpp_seccomp_rules(C.CString(_config.ExePath), C.int(1)))
}
