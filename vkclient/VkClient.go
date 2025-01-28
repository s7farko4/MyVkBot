package vkclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const defaultAPIVersion = "5.199" // Версия API ВК

type VkConfig struct {
	BaseURL           string `json:"base_url"`
	Lang              string `json:"lang"`
	Version           string `json:"version"`
	Token             string `json:"token"`
	ClientID          string `json:"client_id"`
	GroupId           string `json:"group_id"`
	GroupToken        string `json:"group_token"`
	GroupTokenCosplay string `json:"group_token_cosplay"`
	GroupIdCosplay    string `json:"group_id_cosplay"`
	TokenFake         string `json:"token_fake"`
}

type VkResponse struct {
	Response interface{} `json:"response"`
	Error    *VkError    `json:"error,omitempty"`
}
type VkError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

type VkClient struct {
	Config *VkConfig
}

func NewVkClient() (*VkClient, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	config := &VkConfig{}

	err = json.Unmarshal(byteValue, config)
	if err != nil {
		return nil, err
	}

	if config.Version == "" {
		config.Version = defaultAPIVersion
	}

	return &VkClient{Config: config}, nil
}

func (c *VkClient) CallMethod(method string, params map[string]string, token string) (*http.Response, error) {
	url := fmt.Sprintf("%s/method/%s?access_token=%s&v=%s", c.Config.BaseURL, method, token, c.Config.Version)
	for key, value := range params {
		url += "&" + key + "=" + value
	}
	fmt.Println("URL: ", url, "\n")
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var vkResp VkResponse
	err = json.Unmarshal(body, &vkResp)
	if err != nil {
		return nil, err
	}

	if vkResp.Error != nil {
		return nil, errors.New(fmt.Sprint(vkResp.Error))
	}

	return response, nil
}

func (c *VkClient) WallPost(params map[string]string) (*http.Response, error) {

	return c.CallMethod("wall.post", params, c.Config.GroupToken)
}

func (c *VkClient) GetWallUploadServer() (*http.Response, error) {
	params := map[string]string{"group_id": c.Config.GroupId}
	return c.CallMethod("photos.getWallUploadServer", params, c.Config.TokenFake)
}

func (c *VkClient) WallCreateComment(params map[string]string) (*http.Response, error) {
	return c.CallMethod("wall.createComment", params, c.Config.TokenFake)
}

func (c *VkClient) GroupsEditManager(params map[string]string) (*http.Response, error) {
	return c.CallMethod("groups.editManager", params, c.Config.Token)
}

//https:\/\/pu.vk.com\/c857608\/ss2170\/upload.php?act=do_add&mid=198653863&aid=-14&gid=224703507&hash=db77bea2d9af2df64cf49804b89e3b1d&rhash=f03e29eb8d6a1688c09ffc8da018c261&swfupload=1&api=1&wallphoto=1
