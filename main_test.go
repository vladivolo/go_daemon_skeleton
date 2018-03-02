package main

import (
	"github.com/vladivolo/go_daemon_skeleton/service"
	log "github.com/vladivolo/lumber"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	if err := daemon_up("conf/config.yaml"); err != nil {
		t.Error("Start failed! ", err)
	}

	bin, err := GetData("http://" + service.ServiceConf().GetName() + "/ping")
	if err != nil {
		t.Error("Http get failed! ", err)
	}

	if strings.Compare(string(bin), "PONG") != 0 {
		t.Error("Responce failed! ", err)
	}

	time.Sleep(1 * time.Second)
}

func GetData(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		log.Error("GetData: return %s", err)
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}
