package main

import (
	"github.com/calabashdad/utiltools"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"seal/conf"
	"strconv"
	"strings"
)

type hlsServer struct {
}

func (hs *hlsServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}

		gGuards.Done()
	}()

	if "false" == conf.GlobalConfInfo.Hls.Enable {
		log.Println("hls server disabled")
		return
	}
	log.Println("start hls server, listen at :", conf.GlobalConfInfo.Hls.HttpListen)

	http.HandleFunc("/live/", handleLive)

	if err := http.ListenAndServe(":"+conf.GlobalConfInfo.Hls.HttpListen, nil); err != nil {
		log.Println("start hls server failed, err=", err)
	}
}

var crossdomainxml = []byte(
	`<?xml version="1.0" ?><cross-domain-policy>
			<allow-access-from domain="*" />
			<allow-http-request-headers-from domain="*" headers="*"/>
		</cross-domain-policy>`)

func handleLive(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	if path.Base(r.URL.Path) == "crossdomain.xml" {

		w.Header().Set("Content-Type", "application/xml")
		w.Write(crossdomainxml)
		return
	}

	ext := path.Ext(r.URL.Path)
	switch ext {
	case ".m3u8":
		app, m3u8 := parseM3u8File(r.URL.Path)
		m3u8 = conf.GlobalConfInfo.Hls.HlsPath + "/" + app + "/" + m3u8
		if data, err := loadFile(m3u8); nil == err {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Content-Type", "application/x-mpegURL")
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			if _, err = w.Write(data); err != nil {
				log.Println("write m3u8 file err=", err)
			}
		}
	case ".ts":
		app, ts := parseTsFile(r.URL.Path)
		ts = conf.GlobalConfInfo.Hls.HlsPath + "/" + app + "/" + ts
		if data, err := loadFile(ts); nil == err {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "video/mp2ts")
			w.Header().Set("Content-Length", strconv.Itoa(len(data)))
			if _, err = w.Write(data); err != nil {
				log.Println("write ts file err=", err)
			}
		}
	case ".flv":
		u := r.URL.Path

		path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".flv")
		paths := strings.SplitN(path, "/", 2)

		if len(paths) != 2 {
			http.Error(w, "http-flv path error, should be /live/stream.flv", http.StatusBadRequest)
			return
		}
		log.Println("url:", u, "path:", path, "paths:", paths)

		streamName := paths[1]
		w.Header().Set("Access-Control-Allow-Origin", "*")

		httpFlvStreamCycle(streamName, w)
	default:
		log.Println("unknown hls request file, type=", ext)
	}
}

func parseM3u8File(p string) (app string, m3u8File string) {
	if i := strings.Index(p, "/"); i >= 0 {
		if j := strings.LastIndex(p, "/"); j > 0 {
			app = p[i+1 : j]
		}
	}

	if i := strings.LastIndex(p, "/"); i > 0 {
		m3u8File = p[i+1:]
	}

	return
}

func parseTsFile(p string) (app string, tsFile string) {
	if i := strings.Index(p, "/"); i >= 0 {
		if j := strings.LastIndex(p, "/"); j > 0 {
			app = p[i+1 : j]
		}
	}

	if i := strings.LastIndex(p, "/"); i > 0 {
		tsFile = p[i+1:]
	}

	return
}

func loadFile(filename string) (data []byte, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	var f *os.File
	if f, err = os.Open(filename); err != nil {
		log.Println("Open file ", filename, " failed, err is", err)
		return
	}
	defer f.Close()

	if data, err = ioutil.ReadAll(f); err != nil {
		log.Println("read file ", filename, " failed, err is", err)
		return
	}

	return
}

func httpFlvStreamCycle(streamName string, w http.ResponseWriter) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	flvHeader := []byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09}

	var err error
	if _, err = w.Write(flvHeader); err != nil {
		return
	}

}
