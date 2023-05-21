package main

import (
	"log"
	UserServer "tiktok-demo/shared/kitex_gen/UserServer/userservice"
)

func main() {
	svr := UserServer.NewServer(new(UserServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
