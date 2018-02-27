package service

import (
	"flag"
	"fmt"
	log "github.com/vladivolo/lumber"
	"os"
)

// Command line flags
type Flags struct {
	ConfFile string
	LogFile  string
	Version  bool
}

const VersionInfo string = "v0.0.1"

var (
	flags  = Flags{}
	config *ServiceConfig
)

func Initialize(conf_path string) {
	var err error

	flag.StringVar(&flags.ConfFile, "c", conf_path, "path to config file")
	flag.StringVar(&flags.LogFile, "l", "", "path to log file, special value '-' means 'stdout'")
	flag.BoolVar(&flags.Version, "v", false, "print version")
	flag.Parse()

	if flags.Version {
		fmt.Println(VersionInfo)
		os.Exit(0)
	}

	if flags.ConfFile != "" {
		conf_path = flags.ConfFile
	}

	config, err = ConfigInit(conf_path)
	if err != nil {
		log.Error("can't load config: %s", err)
		os.Exit(-11)
	}

	if flags.LogFile != "" {
		config.LogFile = flags.LogFile
	}

	log_level := config.GetLogLevel()
	if log_level == 0 {
		log.Error("unknown log_level")
		os.Exit(-13)
	}

	err = Reopen_logfile(config.GetLogFile(), log_level)
	if err != nil {
		log.Error("can't open logfile: %s", err)
		os.Exit(-14)
	}

}

func CmdlineFlags() *Flags {
	return &flags
}

func ServiceConf() *ServiceConfig {
	return config
}
