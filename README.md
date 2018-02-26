# Seal

seal is rtmp server written by go language, main refer to rtmp server open source https://github.com/ossrs/srs

## Usage
* build

  download https://github.com/calabashdad/seal to ```go path```, run ```go build```
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

## support
  go is cross platform 
* linux
* mac
* windows

## support
* rtmp protocol

## plan to support
* hls
* http-flv
* rtsp