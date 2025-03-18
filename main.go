package main

import (
	"VkBot/GoogleSheets"
	"VkBot/WorkerClient"
	"VkBot/vkclient"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func main() {
	/*
		pClient := parserclient.NewParserClient()
		sources, err := pClient.ParsSources()
		if err != nil {
			log.Fatalf("Ошибка при обработке источников: %v", err)
		}
		fmt.Println(sources.Dirs[0])

		client, err := vkclient.NewVkClient()
		if err != nil {
			panic(err)
		}

		link := "https://t.me/+r5P_voWHdEZkOTFh"
		commentText := "👇Фотки тут👇\n\n" + link + "\n\n" + link + "\n\n👆Нажимай👆"
		messageText := "Очень горячий фотосет получился❤\nСамое горячее в комментариях👇"
		params := map[string]string{
			"commentText": commentText,
			"messageText": messageText,
			"userID":      client.Config.UserID,            //UserID
			"groupToken":  client.Config.GroupTokenCosplay, //GroupTokenCosplay
			"clientID":    client.Config.ClientIdCosplay,   //ClientIdCosplay
			"groupId":     client.Config.GroupIdCosplay,    //GroupIdCosplay
			"token":       client.Config.TokenFake,         //TokenFake
			"ownerToken":  client.Config.Token,             //Token
			//"primary_attachments_mode": "",
		}
		targetTimePost1 := time.Date(2025, 3, 16, 22, 00, 00, 000, time.UTC)

		resp, err := client.TimerPost(targetTimePost1, params, sources.Dirs[0].Photos, sources.Dirs[0].Videos)
		if err != nil {
			panic(err)
		}
		fmt.Println(resp.PostID)
		fmt.Println(resp.PostLink)

		select {}
	*/

	// Создание клиента VK
	vkClient, err := vkclient.NewVkClient()
	if err != nil {
		log.Fatalf("Ошибка создания клиента VK: %v", err)
	}
	googleSheetsClient, err := GoogleSheets.NewClient()
	if err != nil {
		panic(err)
	}
	worker := WorkerClient.NewWorker(*vkClient, *googleSheetsClient)

	err = worker.StartWorker()
	if err != nil {
		panic(err)
	}
	select {}

}
