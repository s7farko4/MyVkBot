package main

import (
	"VkBot/ParserClient"
	"VkBot/vkclient"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Драйвер для SQLite
	"log"
	"time"
)

func main() {
	pClient := parserclient.NewParserClient()
	sources, err := pClient.ParsSources()
	if err != nil {
		log.Fatalf("Ошибка при обработке источников: %v", err)
	}
	fmt.Println(sources)

	client, err := vkclient.NewVkClient()
	if err != nil {
		panic(err)
	}
	commentText := "Больше фото тут 👉 https://t.me/+mn_aRMGgiEI5ZWVk"
	messageText := "Подписчица поделилась откровенными фото💦\nПродолжение в комментариях👇"

	params := map[string]string{
		"commentText": commentText,
		"messageText": messageText,
		"userID":      client.Config.UserID,
		"groupToken":  client.Config.GroupToken,
		"clientID":    client.Config.ClientID,
		"groupId":     client.Config.GroupId,
		"token":       client.Config.TokenFake,
		"ownerToken":  client.Config.Token,
	}
	// Устанавливаем конкретную дату и время
	targetTime := time.Date(2025, 1, 30, 00, 00, 00, 000, time.UTC)

	resp, err := client.TimerPost(targetTime, params, sources.Dirs[0].Photos)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	select {}

}
