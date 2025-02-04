package main

import (
	parserclient "VkBot/ParserClient"
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
	commentText := "Больше фото тут 👉 https://t.me/+_CKpLbxW5QtkMzky"
	messageText := "🍭 Mилaшечка пoдписчицa cкинула свoи фoтки на oцeнкy 🥵 Её кaнaльчик в кoммeнтаpиях 😳"

	params := map[string]string{
		"commentText": commentText,
		"messageText": messageText,
		"userID":      client.Config.UserID,
		"groupToken":  client.Config.GroupToken,
		"clientID":    client.Config.ClientID,
		"groupId":     client.Config.GroupId,
		"token":       client.Config.TokenFake,
		"ownerToken":  client.Config.Token,
		//"primary_attachments_mode": "",
	}
	// Устанавливаем конкретную дату и время
	targetTime := time.Date(2025, 1, 30, 22, 00, 00, 000, time.UTC)

	resp, err := client.TimerPost(targetTime, params, sources.Dirs[0].Photos, sources.Dirs[0].Videos)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	select {}
	/*
		resp, uploadUrl, err := client.VideoSave(params["token"])
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
		fmt.Println(uploadUrl)

		res, err := vkclient.UploadVideo("sources/Beautiful Cosplay/1.mp4", uploadUrl)
		fmt.Println(int64((res["video_id"]).(float64)))
	*/

}
