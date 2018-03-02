package main

import (
	"github.com/vladivolo/go_daemon_skeleton/service"
	log "github.com/vladivolo/lumber"
	"os"
)

func main() {
	if err := daemon_up("/default/path/to/config.yaml"); err != nil {
		os.Exit(-253)
	}

	service.Wait_for_signals(sigaction__graceful_shutdown, sigaction__reopen_logs, sigaction__reload_config)
}

func daemon_up(config string) error {
	service.Initialize(config, VersionInfo.Version)

	StartHttpServer(service.ServiceConf().GetListen(), service.ServiceConf().GetHttpWorkersCount())

	log.Info("Start Daemon: %v", service.ServiceConf())

	return nil
}
