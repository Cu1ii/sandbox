package sandbox

func _cCppSeccompRules(_config *config, allow_wirte_file bool) int {
	return 0
}

func cCppSeccompRules(_config *config) int {
	return _cCppSeccompRules(_config, false)
}

func cCppFileIoSeccompRules(_config *config) int {
	return _cCppSeccompRules(_config, true)
}
