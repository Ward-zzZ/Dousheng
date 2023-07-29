package main

import (
	"bytes"
	"os"
	"log"
	"net/http"
	"strings"
)

func main() {
	// 读取文件内容
	data, err := os.ReadFile("data/configs/consul/config.yaml")
	if err != nil {
		panic(err)
	}

	// 将字节数组转换为字符串
	content := string(data)

	// 使用字符串分割函数将字符串分割成多个内容块
	blocks := strings.Split(content, "\n\n")

	kv := map[string]string{
		"tiktok/api_srv":      blocks[0],
		"tiktok/comment_srv":  blocks[1],
		"tiktok/favorite_srv": blocks[2],
		"tiktok/message_srv":  blocks[3],
		"tiktok/relation_srv": blocks[4],
		"tiktok/user_srv":     blocks[5],
		"tiktok/video_srv":    blocks[6],
		//添加更多键值对...
	}

	client := &http.Client{}

	for key, value := range kv {
		request, err := http.NewRequest("PUT", "http://localhost:8500/v1/kv/"+key, bytes.NewBufferString(value))
		if err != nil {
			log.Fatalf("Unable to create request: %v", err)
		}

		response, err := client.Do(request)
		if err != nil {
			log.Fatalf("Unable to send request: %v", err)
		} else {
			defer response.Body.Close()
			log.Printf("Response status code for key '%s': %d", key, response.StatusCode)
		}
	}
}
