package main

import (
	"encoding/binary"
	"github.com/calabashdad/utiltools"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"seal/conf"
	"seal/rtmp/co"
	"strconv"
	"strings"
	"time"
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

		key := paths[0] + "/" + paths[1]
		w.Header().Set("Access-Control-Allow-Origin", "*")

		httpFlvStreamCycle(key, w)
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

func httpFlvStreamCycle(key string, w http.ResponseWriter) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(utiltools.PanicTrace())
		}
	}()

	var err error

	source := co.GlobalSources.FindSourceToPlay(key)
	if nil == source {
		log.Printf("httpFlvStreamCycle, stream=%s can not play because has not published\n", key)
		http.Error(w, "this stream has not published", http.StatusBadRequest)
		return
	}

	consumer := co.NewConsumer()
	source.CreateConsumer(consumer)

	if source.Atc && !source.GopCache.Empty() {
		if nil != source.CacheMetaData {
			source.CacheMetaData.Header.Timestamp = source.GopCache.StartTime()
		}
		if nil != source.CacheVideoSequenceHeader {
			source.CacheVideoSequenceHeader.Header.Timestamp = source.GopCache.StartTime()
		}
		if nil != source.CacheAudioSequenceHeader {
			source.CacheAudioSequenceHeader.Header.Timestamp = source.GopCache.StartTime()
		}
	}

	//cache meta data
	if nil != source.CacheMetaData {
		consumer.Enquene(source.CacheMetaData, source.Atc, source.SampleRate, source.FrameRate, source.TimeJitter)
	}

	//cache video data
	if nil != source.CacheVideoSequenceHeader {
		consumer.Enquene(source.CacheVideoSequenceHeader, source.Atc, source.SampleRate, source.FrameRate, source.TimeJitter)
	}

	//cache audio data
	if nil != source.CacheAudioSequenceHeader {
		consumer.Enquene(source.CacheAudioSequenceHeader, source.Atc, source.SampleRate, source.FrameRate, source.TimeJitter)
	}

	//dump gop cache to client.
	source.GopCache.Dump(consumer, source.Atc, source.SampleRate, source.FrameRate, source.TimeJitter)

	log.Println("httpFlvStreamCycle now playing, key=", key)

	// send flv header
	flvHeader := []byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09}
	if _, err = w.Write(flvHeader); err != nil {
		log.Println("httpFlvStreamCycle send flv header to remote success.")
		return
	}

	timeLast := time.Now().Unix()

	var previousTagLen uint32
	for {
		msg := consumer.Dump()
		if nil == msg {
			// wait and try again.

			timeCurrent := time.Now().Unix()
			if timeCurrent-timeLast > 30 {
				log.Println("httpFlvStreamCycle time out > 30, break. key=", key)
				break
			}

			continue
		} else {

			timeLast = time.Now().Unix()

			// previous tag len c4B. 11 + payload data size
			// type 1B
			// data size 3B
			// timestamp 3B
			// timestampEx 1B
			// streamID 3B always is 0
			// total is 4 + 1 +3 + 3 + 1 + 3 = 15B
			var tagHeader [15]uint8
			var offset uint32

			// previous tag len
			binary.BigEndian.PutUint32(tagHeader[offset:], previousTagLen)
			offset += 4

			// type
			tagHeader[offset] = msg.Header.MessageType
			offset++

			// payload data size
			var sizebuf [4]uint8
			binary.BigEndian.PutUint32(sizebuf[:], msg.Header.PayloadLength)
			copy(tagHeader[offset:], sizebuf[1:])
			offset += 3

			// timestamp
			var timebuf [4]uint8
			binary.BigEndian.PutUint32(timebuf[:], uint32(msg.Header.Timestamp))
			copy(tagHeader[offset:], timebuf[1:])
			offset += 3

			// timestamp ex, generally not used
			tagHeader[offset] = 0
			offset++

			// stream id
			tagHeader[offset] = 0
			offset++
			tagHeader[offset] = 0
			offset++
			tagHeader[offset] = 0
			offset++

			if _, err = w.Write(tagHeader[:]); err != nil {
				log.Println("httpFlvStreamCycle: playing... send tag header to remote failed.err=", err)
				break
			}

			if _, err = w.Write(msg.Payload.Payload); err != nil {
				log.Println("httpFlvStreamCycle: playing... send tag payload to remote failed.err=", err)
				break
			}

			previousTagLen = 11 + msg.Header.PayloadLength
		}
	}

	source.DestroyConsumer(consumer)
	log.Printf("httpFlvStreamCycle: playing over, key=%s, consumer has destroyed\n", key)

}
