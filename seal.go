package main

import (
	"UtilsTools/identify_panic"
	"flag"
	"log"
	"os"
	"sync"
	"time"
)

const SEAL_VERSION = "seal version: 1.0.0"

var (
	conf_file    = flag.String("c", "./seal.yaml", "configure filename")
	show_version = flag.Bool("v", false, "show version of seal")
)

var (
	g_wg sync.WaitGroup
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
		return
	}

	g_wg.Add(1)
	StartRtmpServer()

	g_wg.Wait()
	log.Println("seal quit gracefully.")
}
