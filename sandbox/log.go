package sandbox

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

const (
	LogLevelFatal = iota
	LogLevelWarning
	LogLevelInfo
	LogLevelDebug
)

var (
	logLevelNote = [5]string{"FATAL", "WARNING", "INFO", "DEBUG"}
)

func LogDebug(logFp *os.File, message string) {
	logWrite(LogLevelDebug, logFp, message)
}

func LogInfo(logFp *os.File, message string) {
	logWrite(LogLevelInfo, logFp, message)
}

func LogWarning(logFp *os.File, message string) {
	logWrite(LogLevelWarning, logFp, message)
}

func LogFatal(logFp *os.File, message string) {
	logWrite(LogLevelFatal, logFp, message)
}

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func LogOpen(filename string) *os.File {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Println("can not open log file " + filename)
		return nil
	}
	return file
}

func CloseFile(file *os.File) {
	file.Close()
}

func logWrite(level int, logFp *os.File, message string) {
	if logFp == nil {
		log.Println("can not open log file " + logFp.Name())
		return
	}
	if err := syscall.Flock(int(logFp.Fd()), syscall.LOCK_EX); err == nil {
		log.SetPrefix("[" + logLevelNote[level] + "] ")
		log.SetOutput(logFp)
		log.Println(message)
		_ = syscall.Flock(int(logFp.Fd()), syscall.LOCK_UN)
	} else {
		fmt.Fprintln(os.Stderr, "flock error")
	}

}
