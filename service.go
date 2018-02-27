package main

import (
	"github.com/vladivolo/go_daemon_skeleton/service"
	log "github.com/vladivolo/lumber"
	"os"
)

func sigaction__graceful_shutdown(sig os.Signal) {
	log.Info("Got %s, shutting down", service.SignalName(sig))

	CloseHttpServer()
}

func sigaction__reopen_logs(sig os.Signal) {
	log.Info("Got %s, reopening logfile: %s", service.SignalName(sig), service.LogPath())
	service.Reopen_logfile(service.LogPath(), service.LogLevel())
}

func sigaction__reload_config(sig os.Signal) {
}
