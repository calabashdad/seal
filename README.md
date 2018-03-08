# Seal

seal is rtmp server written by go language, main refer to rtmp server open source https://github.com/ossrs/srs

## Usage
* build

  download https://github.com/calabashdad/seal to ```go path```, run ```go build```

  you can also use cross platform build, like build a linux version if you are on mac, run ```cross_platform_linux```

* run console mode

  ```./seal -c seal.yaml```
* run daemon mode

  ```nohup ./seal -c seal.yaml &```
* mock stream publish
  
  <pre><code>for((;;)); do \
        ffmpeg -re -i lindan.flv \
        -vcodec copy -acodec copy \
        -f flv -y rtmp://127.0.0.1/live/test; \
	    sleep 3       
  done</code></pre> 

* use vlc play 
```rtmp://127.0.0.1/live/test```

## platform
  go is cross platform 
* linux
* mac
* windows

## support
* rtmp protocol (h264)

## plan to support
* hls 
* rtsp
* http-flv
* h265
* transcode(audio to aac)
* http stats query
* video on demand
* video encry
* auth token dynamicly
* mini rtmp server in embed device