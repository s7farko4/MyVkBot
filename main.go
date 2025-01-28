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
	filePath := "C:/Users/s7far/GolandProjects/VkBot/sorces/a.jpg"
	commentText := "Больше фото тут 👉 https://t.me/+E-DuB-Axd6RhMmFh"
	messageText := "Это рекламная запись, переходите по ссылке в комментариях"
	groupId := client.Config.GroupId
	userId := client.Config.UserID
	clientId := client.Config.ClientID

	//Устанавливает пользователя редактором
	paramsEditManager := map[string]string{
		"group_id": groupId,
		"user_id":  userId,
		"role":     "editor",
	}
	resp, err := client.GroupsEditManager(paramsEditManager)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	//Оставляет запись на стене сообщества
	att, err := client.GetAttachments(filePath)
	if err != nil {
		panic(err)
	}
	paramsWallPoast := map[string]string{
		"owner_id":     clientId,
		"from_group":   groupId,
		"attachments":  att,
		"publish_date": "",
	}
	resp, postID, err := client.WallPost(paramsWallPoast, messageText)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	//пишет комментарий под постом от имени сообщества

	paramsCreateComment := map[string]string{
		"owner_id":   clientId,
		"post_id":    postID,
		"from_group": groupId,
	}

	resp, err = client.WallCreateComment(paramsCreateComment, commentText)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	//удаляет все роли у выбранного пользователя
	paramsEditManager = map[string]string{
		"group_id": groupId,
		"user_id":  userId,
		"role":     "",
	}
	resp, err = client.GroupsEditManager(paramsEditManager)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

}
