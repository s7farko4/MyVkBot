package main

import (
	"VkBot/vkclient"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Драйвер для SQLite
	"time"
)

func main() {
	client, err := vkclient.NewVkClient()
	if err != nil {
		panic(err)
	}
	imagePaths := make([]string, 2)
	imagePaths[0] = "sources/1_2.jpg"
	imagePaths[1] = "sources/1_1.jpg"
	commentText := "Больше фото тут 👉 https://t.me/+mn_aRMGgiEI5ZWVk"
	messageText := "Подписчица поделилась откровенными фото💦\nПродолжение в комментариях👇"

	params := map[string]string{
		"commentText": commentText,
		"messageText": messageText,
		"userID":      client.Config.UserIDCosplay,
		"groupToken":  client.Config.GroupTokenCosplay,
		"clientID":    client.Config.ClientIdCosplay,
		"groupId":     client.Config.GroupIdCosplay,
		"token":       client.Config.TokenFakeCosplay,
		"ownerToken":  client.Config.Token,
	}
	// Устанавливаем конкретную дату и время
	targetTime := time.Date(2025, 1, 30, 00, 00, 00, 000, time.UTC)

	resp, err := client.TimerPost(targetTime, params, imagePaths)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	select {}

}
