package conf

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/calabashdad/utiltools"
	"github.com/yaml"
)

// GlobalConfInfo global config info
var GlobalConfInfo confInfo

type systemConfInfo struct {
	CPUNums uint32 `yaml:"cpuNums"`
}

type rtmpConfInfo struct {
	Listen            string `yaml:"listen"`
	TimeOut           uint32 `yaml:"timeout"`
	ChunkSize         uint32 `yaml:"chunkSize"`
	Atc               bool   `yaml:"atc"`
	AtcAuto           bool   `yaml:"atcAuto"`
	TimeJitter        uint32 `yaml:"timeJitter"`
	ConsumerQueueSize uint32 `yaml:"consumerQueueSize"`
}

type hlsConfInfo struct {
	Enable      string `yaml:"enable"`
	HlsFragment int    `yaml:"hlsFragment"`
	HlsWindow   int    `yaml:"hlsWindow"`
	HlsPath     string `yaml:"hlsPath"`
}

type confInfo struct {
	System systemConfInfo `yaml:"system"`
	Rtmp   rtmpConfInfo   `yaml:"rtmp"`
	Hls    hlsConfInfo    `yaml:"hls"`
}

func (t *confInfo) Loads(c string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}

	}()

	var f *os.File
	if f, err = os.Open(c); err != nil {
		log.Println("Open config failed, err is", err)
		return err
	}
	defer f.Close()

	var data []byte
	if data, err = ioutil.ReadAll(f); err != nil {
		log.Println("config file loads failed, ", err.Error())
		return err
	}

	if err = yaml.Unmarshal(data, t); err != nil {
		log.Println("error:", err.Error())
		return err
	}

	return nil
}
