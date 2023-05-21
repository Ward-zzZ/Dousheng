package main

import (
	"log"
	FavoriteServer "tiktok-demo/shared/kitex_gen/FavoriteServer/favoriteservice"
)

func main() {
	svr := FavoriteServer.NewServer(new(FavoriteServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
