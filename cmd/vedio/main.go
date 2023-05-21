package main

import (
	"log"
	VideoServer "tiktok-demo/shared/kitex_gen/VideoServer/videosrv"
)

func main() {
	svr := VideoServer.NewServer(new(VideoSrvImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
