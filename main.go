package main

import (
	"VkBot/vkclient"
	"fmt"
	"net/url"
)

func main() {
	commentText := "Больше фото тут 👉 https://t.me/+E-DuB-Axd6RhMmFh"
	escapedComment := url.QueryEscape(commentText)
	messageText := "Это рекламная запись, переходите по ссылке в комментариях"
	escapedMessage := url.QueryEscape(messageText)
	client, err := vkclient.NewVkClient()
	if err != nil {
		panic(err)
	}

	//Устанавливает пользователя редактором
	paramsEditManager := map[string]string{
		"group_id": client.Config.GroupId,
		"user_id":  client.Config.UserID,
		"role":     "editor",
	}
	resp, err := client.GroupsEditManager(paramsEditManager)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	//Оставляет запись на стене сообщества
	paramsWallPoast := map[string]string{
		"owner_id":     client.Config.ClientID,
		"from_group":   client.Config.GroupId,
		"message":      escapedMessage,
		"attachments":  "",
		"publish_date": "",
	}
	resp, postID, err := client.WallPost(paramsWallPoast)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	//пишет комментарий под постом от имени сообщества

	paramsCreateComment := map[string]string{
		"owner_id":   client.Config.ClientID,
		"post_id":    postID,
		"from_group": client.Config.GroupId,
		"message":    escapedComment,
	}

	resp, err = client.WallCreateComment(paramsCreateComment)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

	//удаляет все роли у выбранного пользователя
	paramsEditManager = map[string]string{
		"group_id": client.Config.GroupId,
		"user_id":  client.Config.UserID,
		"role":     "",
	}
	resp, err = client.GroupsEditManager(paramsEditManager)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}

}
