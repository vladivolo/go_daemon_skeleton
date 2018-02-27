package service

import (
	log "github.com/vladivolo/lumber"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type StorageData struct {
	Path   string `json:"path" yaml:"path"`
	Enable bool   `json:"enable" yaml:"enable"`
}

type ServiceConfig struct {
	Name              string        `json:"name" yaml:"name"`
	Listen            string        `json:"listen" yaml:"listen"`
	Http_wokers_count int           `json:"http_workers_count" yaml:"http_workers_count"`
	Pgx               string        `json:"pgx" yaml:"pgx"`
	Output_queue      string        `json:"output_queue" yaml:"output_queue"`
	Input_queue       string        `json:"input_queue" yaml:"input_queue"`
	LogFile           string        `json:"logfile" yaml:"logfile"`
	LogLevel          string        `json:"loglevel" yaml:"loglevel"`
	Storage           []StorageData `json:"storage" yaml:"storage"`
}

func ConfigInit(filePath string) (*ServiceConfig, error) {
	var Conf = ServiceConfig{}

	if _, err := os.Stat(filePath); err != nil {
		log.Error("File (%s) (%s)", filePath, err)
		return nil, err
	}

	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error("Could not read file %s error %s", filePath, err)
		return nil, err
	}

	if err := yaml.Unmarshal(buf, &Conf); err != nil {
		log.Error("Syntax error in %s: %s", filePath, err)
		return nil, err
	}

	return &Conf, nil
}

func (s *ServiceConfig) GetName() string {
	return s.Name
}

func (s *ServiceConfig) GetListen() string {
	return s.Listen
}

func (s *ServiceConfig) GetHttpWorkersCount() int {
	if s.Http_wokers_count > 0 {
		return s.Http_wokers_count
	}
	return 1
}

func (s *ServiceConfig) GetPgx() string {
	return s.Pgx
}

func (s *ServiceConfig) GetInputQueue() string {
	return s.Input_queue
}

func (s *ServiceConfig) GetOutputQueue() string {
	return s.Output_queue
}

func (s *ServiceConfig) GetStorage() (ss []string) {
	for _, x := range s.Storage {
		if x.Enable {
			ss = append(ss, x.Path)
		}
	}

	return
}

func (s *ServiceConfig) GetLogFile() string {
	return s.LogFile
}

func (s *ServiceConfig) GetLogLevel() int {
	return log.LvlInt(s.LogLevel)
}
