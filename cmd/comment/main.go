package main

import (
	"log"
	CommentServer "tiktok-demo/shared/kitex_gen/CommentServer/commentservice"
)

func main() {
	svr := CommentServer.NewServer(new(CommentServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
