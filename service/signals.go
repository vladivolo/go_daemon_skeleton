package service

import (
	"fmt"
	log "github.com/vladivolo/lumber"
	"os"
	"os/signal"
	"syscall"
)

func SignalName(sig os.Signal) string {
	switch sig {
	case syscall.SIGABRT:
		return "SIGABRT"
	case syscall.SIGALRM:
		return "SIGALRM"
	case syscall.SIGBUS:
		return "SIGBUS"
	case syscall.SIGCHLD:
		return "SIGCHLD"
	case syscall.SIGCONT:
		return "SIGCONT"
	case syscall.SIGFPE:
		return "SIGFPE"
	case syscall.SIGHUP:
		return "SIGHUP"
	case syscall.SIGILL:
		return "SIGILL"
	case syscall.SIGINT:
		return "SIGINT"
	case syscall.SIGIO:
		return "SIGIO"
	case syscall.SIGKILL:
		return "SIGKILL"
	case syscall.SIGPIPE:
		return "SIGPIPE"
	case syscall.SIGPROF:
		return "SIGPROF"
	case syscall.SIGPWR:
		return "SIGPWR"
	case syscall.SIGQUIT:
		return "SIGQUIT"
	case syscall.SIGSEGV:
		return "SIGSEGV"
	case syscall.SIGSTKFLT:
		return "SIGSTKFLT"
	case syscall.SIGSTOP:
		return "SIGSTOP"
	case syscall.SIGSYS:
		return "SIGSYS"
	case syscall.SIGTERM:
		return "SIGTERM"
	case syscall.SIGTRAP:
		return "SIGTRAP"
	case syscall.SIGTSTP:
		return "SIGTSTP"
	case syscall.SIGTTIN:
		return "SIGTTIN"
	case syscall.SIGTTOU:
		return "SIGTTOU"
	case syscall.SIGURG:
		return "SIGURG"
	case syscall.SIGUSR1:
		return "SIGUSR1"
	case syscall.SIGUSR2:
		return "SIGUSR2"
	case syscall.SIGVTALRM:
		return "SIGVTALRM"
	case syscall.SIGWINCH:
		return "SIGWINCH"
	case syscall.SIGXCPU:
		return "SIGXCPU"
	case syscall.SIGXFSZ:
		return "SIGXFSZ"
	default:
		return fmt.Sprintf("SIGNAL (%s)", sig)
	}
}

func Wait_for_signals(shutdown func(sig os.Signal), reopen_log func(sig os.Signal), reload_conf func(sig os.Signal)) {
	c := make(chan os.Signal, 5)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT)

	for {
		sig := <-c

		log.Info("Got signal: %s", SignalName(sig))

		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			shutdown(sig)
			return
		case syscall.SIGUSR1:
			reopen_log(sig)
		case syscall.SIGUSR2:
			reload_conf(sig)
		}
	}
}
