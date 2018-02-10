package main

import (
	"flag"
	"log"
	"os"
	"seal/conf"
	"sync"
	"time"
	"utiltools"
)

const SealVersion = "seal: 1.0.0"

var (
	configFile  = flag.String("c", "./seal.yaml", "configure filename")
	showVersion = flag.Bool("v", false, "show version of seal")
)

var (
	gWgServers sync.WaitGroup
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate |log.Lmicroseconds)
	flag.Parse()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
			time.Sleep(1 * time.Second)
		}
	}()

	if len(os.Args) < 2 {
		log.Println("Show usage : ./seal --help.")
		return
	}

	if *showVersion {
		log.Println(SealVersion)
		return
	}

	err := conf.GlobalConfInfo.Loads(*configFile)
	if err != nil {
		log.Println("conf loads failed.err=", err)

		//load conf file failed. use default config.
		conf.GlobalConfInfo.Default()
	} else {
		log.Println("load conf file success, conf=", conf.GlobalConfInfo)
	}

	gWgServers.Add(1)
	if true {
		rtmpServer := RtmpServer{}

		rtmpServer.Start()
	}

	gWgServers.Wait()
	log.Println("seal quit gracefully.")
}
