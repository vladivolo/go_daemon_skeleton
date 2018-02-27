package service

import (
	log "github.com/vladivolo/lumber"
	"os"
)

var (
	logPath    = "-"
	logLevel   = log.ERROR
	logBufsize = 1024
)

func LogPath() string {
	return logPath
}

func LogLevel() int {
	return logLevel
}

func Reopen_logfile(path string, level int) (err error) {
	var l log.Logger

	if path == "" || path == "-" {
		l = log.NewBasicLogger(os.Stderr, level)
	} else {
		l, err = log.NewFileLogger(path, level, log.APPEND, 0, 0, logBufsize)
		if err != nil {
			return
		}

		l.Info("reopen_log_file: new log opened with path: %s\n", path)
	}

	log.SetLogger(l)

	logPath = path
	logLevel = level

	return nil
}
