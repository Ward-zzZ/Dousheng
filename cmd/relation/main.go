package main

import (
	"log"
	RelationServer "tiktok-demo/shared/kitex_gen/RelationServer/relationservice"
)

func main() {
	svr := RelationServer.NewServer(new(RelationServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
