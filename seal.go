package main

import (
	"UtilsTools/identify_panic"
	"flag"
	"log"
	"os"
	"seal/conf"
	"sync"
	"time"
)

const SEAL_VERSION = "seal: 1.0.0"

var (
	conf_file    = flag.String("c", "./seal.yaml", "configure filename")
	show_version = flag.Bool("v", false, "show version of seal")
)

var (
	gWGServers sync.WaitGroup
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, ",panic at ", identify_panic.IdentifyPanic())
			time.Sleep(1 * time.Second)
		}
	}()

	if len(os.Args) < 2 {
		log.Println("Show usage : ./seal --help.")
		return
	}

	if *show_version {
		log.Println(SEAL_VERSION)
		return
	}

	err := conf.GlobalConfInfo.Loads(*conf_file)
	if err != nil {
		log.Println("conf loads failed.err=", err)

		//load conf file failed. use default config.
		conf.GlobalConfInfo.Default()
	} else {
		log.Println("load conf file success, conf=", conf.GlobalConfInfo)
	}

	gWGServers.Add(1)
	if true {
		rtmp_server := RtmpServer{}

		rtmp_server.Start()
	}

	gWGServers.Wait()
	log.Println("seal quit gracefully.")
}
