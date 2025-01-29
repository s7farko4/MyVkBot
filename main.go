package main

import (
	"VkBot/vkclient"
	"fmt"
)

func main() {

	client, err := vkclient.NewVkClient()
	if err != nil {
		panic(err)
	}
	imagePaths := make([]string, 2)
	imagePaths[0] = "sorces/1_1.jpg"
	imagePaths[1] = "sorces/1_2.jpg"
	commentText := "Больше фото тут 👉 https://t.me/+E-DuB-Axd6RhMmFh"
	messageText := "Это рекламная запись, переходите по ссылке в комментариях"

	params := map[string]string{
		"commentText": commentText,
		"messageText": messageText,
	}
	resp, err := client.PostWithOpt(imagePaths, params, true, true, true, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)

}
