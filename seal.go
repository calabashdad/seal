package main

import (
	"UtilsTools/identify_panic"
	"flag"
	"log"
	"os"
	"seal/seal_conf"
	"seal/seal_rtmp_server"
	"sync"
	"time"
)

const SEAL_VERSION = "seal: 1.0.0"

var (
	conf_file    = flag.String("c", "./seal.yaml", "configure filename")
	show_version = flag.Bool("v", false, "show version of seal")
)

var (
	g_conf_info seal_conf.ConfInfo
	g_wg        sync.WaitGroup
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err, identify_panic.IdentifyPanic())
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

	err := g_conf_info.Loads(*conf_file)
	if err != nil {
		log.Println("conf loads failed.err=", err)

		//load conf file failed. use default config.
		g_conf_info.Default()
	} else {
		log.Println("load conf file success, conf=", g_conf_info)
	}

	g_wg.Add(1)
	if true {
		rtmpServer := seal_rtmp_server.RtmpServer{
			Conf: &g_conf_info.Rtmp,
			Wg:   &g_wg,
		}
		rtmpServer.Start()
	}

	g_wg.Wait()
	log.Println("seal quit gracefully.")
}
