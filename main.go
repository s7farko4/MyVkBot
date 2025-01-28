package main

import (
	"VkBot/vkclient"
	"fmt"
	"net/url"
)

func main() {
	postID := "116104"
	commentText := "Больше фото тут 👉 https://t.me/+E-DuB-Axd6RhMmFh"
	client, err := vkclient.NewVkClient()
	if err != nil {
		panic(err)
	}
	/*
		escapedMessage := url.QueryEscape(commentText)

		params := map[string]string{
			"owner_id":   client.Config.GroupId,
			"post_id":    postID,
			"from_group": "0",
			"message":    escapedMessage,
		}

		resp, err := client.WallCreateComment(params)
		if err != nil {
			fmt.Println(resp)
			panic(err)
		}
		fmt.Println(resp.Status)*/
	paramsEditManager := map[string]string{
		"group_id": client.Config.GroupIdCosplay,
		"user_id":  "19258661",
		"role":     "editor",
	}
	resp, err := client.GroupsEditManager(paramsEditManager)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	fmt.Println(resp.Status)

	escapedMessage := url.QueryEscape(commentText)

	paramsCreateComment := map[string]string{
		"owner_id":   "-159445401",
		"post_id":    postID,
		"from_group": "159445401",
		"message":    escapedMessage,
	}

	resp, err = client.WallCreateComment(paramsCreateComment)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	fmt.Println(resp.Status)

	paramsEditManager = map[string]string{
		"group_id": client.Config.GroupIdCosplay,
		"user_id":  "19258661",
		"role":     "",
	}
	resp, err = client.GroupsEditManager(paramsEditManager)
	if err != nil {
		fmt.Println(resp)
		panic(err)
	}
	fmt.Println(resp.Status)
}
