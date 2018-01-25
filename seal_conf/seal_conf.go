package seal_conf

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/yaml"
)

type ConfInfoRtmp struct {
	Listen    string `yaml:"listen"`
	TimeOut   uint32 `yaml:"timeout"`
	ChunkSize uint32 `yaml:"chunk_size"`
}

type ConfInfo struct {
	Rtmp ConfInfoRtmp `yaml:"rtmp"`
}

func (t *ConfInfo) Loads(c string) (err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}

	}()

	var f *os.File
	if f, err = os.Open(c); err != nil {
		log.Println("Open config failed, err is", err)
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("config file loads failed, ", err.Error())
		return err
	}

	err = yaml.Unmarshal(data, t)
	if err != nil {
		log.Println("error:", err.Error())
		return err
	}

	return nil
}

func (t *ConfInfo) Default() {
	t.Rtmp.Listen = "1935"
	t.Rtmp.TimeOut = 30
}
