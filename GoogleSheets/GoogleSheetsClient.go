package vkclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
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
	TokenFake         string `json:"token_fake"`
	UserID            string `json:"user_id"`
	GroupTokenCosplay string `json:"group_token_cosplay"`
	GroupIdCosplay    string `json:"group_id_cosplay"`
	ClientIdCosplay   string `json:"client_id_cosplay"`
	TokenFakeCosplay  string `json:"token_fake_cosplay"`
	UserIDCosplay     string `json:"user_id_cosplay"`
}

type VkResponse struct {
	Response interface{} `json:"response"`
	PostID   string
	PostLink string
	Error    *VkError `json:"error,omitempty"`
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

func (c *VkClient) CallMethod(method string, params map[string]string, token string) (VkResponse, error) {
	url := fmt.Sprintf("%s/method/%s?access_token=%s&v=%s", c.Config.BaseURL, method, token, c.Config.Version)
	for key, value := range params {
		url += "&" + key + "=" + value
	}
	fmt.Println("URL  ", method, " :", url, "\n")
	response, err := http.Get(url)
	if err != nil {
		return VkResponse{}, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return VkResponse{}, err
	}

	var vkResp VkResponse
	err = json.Unmarshal(body, &vkResp)
	if err != nil {
		return VkResponse{}, err
	}

	if vkResp.Error != nil {
		return VkResponse{}, errors.New(fmt.Sprint(vkResp.Error))
	}

	return vkResp, nil
}

func (c *VkClient) WallPost(params map[string]string, message string, token string) (VkResponse, string, string, error) {
	escapedMessage := url.QueryEscape(message)
	params["message"] = escapedMessage

	resp, err := c.CallMethod("wall.post", params, token)
	if err != nil {
		return VkResponse{}, "", "", err
	}
	// Приведение значения post_id к типу float64
	postIdFloat := resp.Response.(map[string]interface{})["post_id"].(float64)

	// Преобразование float64 в целое число и затем в строку
	postId := strconv.FormatInt(int64(postIdFloat), 10)
	link := "https://vk.com/wall" + params["owner_id"] + "_" + postId

	return resp, postId, link, nil
}

func (c *VkClient) WallCreateComment(params map[string]string, message string, token string) (VkResponse, error) {
	escapedMessage := url.QueryEscape(message)
	params["message"] = escapedMessage
	return c.CallMethod("wall.createComment", params, token)
}

func (c *VkClient) GetById(id string, token string) (int, error) {
	params := make(map[string]string)
	params["posts"] = id
	resp, err := c.CallMethod("wall.getById", params, token)
	if err != nil {
		return -1, err
	}

	responseArray := resp.Response.([]interface{})

	if len(responseArray) > 0 {
		firstItem := responseArray[0].(map[string]interface{}) // Преобразование интерфейса в map[string]interface{}
		views, ok := firstItem["views"].(map[string]interface{})
		if !ok {
			log.Fatalf("Не удалось найти поле 'views' в первом элементе.")
		}
		count, ok := views["count"].(float64)
		if !ok {
			log.Fatalf("Не удалось найти поле 'count' в 'views'.")
		}
		return int(count), nil
	} else {
		log.Fatal("Массив response пуст.")
		return -1, err
	}
}

func (c *VkClient) GroupsEditManager(params map[string]string, token string) (VkResponse, error) {
	return c.CallMethod("groups.editManager", params, token)
}

func (c *VkClient) GetWallUploadServer(groupID string, token string) (VkResponse, string, error) {
	params := map[string]string{"group_id": groupID}
	resp, err := c.CallMethod("photos.getWallUploadServer", params, token)
	if err != nil {
		return VkResponse{}, "", err
	}
	// Приведение значения post_id к типу float64
	uploadUrl := resp.Response.(map[string]interface{})["upload_url"].(string)

	return resp, uploadUrl, nil
}

func UploadPhoto(filePath, uploadURL string) (map[string]interface{}, error) {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("photo", file.Name())
	if err != nil {
		fmt.Println("Ошибка создания части multipart:", err)
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Ошибка копирования файла:", err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Ошибка закрытия multipart:", err)
		return nil, err
	}

	request, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return nil, err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	return result, nil
}

func (c *VkClient) PhotosSaveWallPhoto(postServerResp map[string]interface{}, token string, groupID string) (VkResponse, error) {

	params := map[string]string{
		"server":   strconv.FormatFloat(postServerResp["server"].(float64), 'f', -1, 64),
		"hash":     postServerResp["hash"].(string),
		"v":        "5.199",
		"photo":    postServerResp["photo"].(string),
		"group_id": groupID,
	}
	return c.CallMethod("photos.saveWallPhoto", params, token)
}

func (c *VkClient) VideoSave(token string) (VkResponse, string, error) {
	params := map[string]string{}
	resp, err := c.CallMethod("video.save", params, token)

	if err != nil {
		return VkResponse{}, "", err
	}
	uploadUrl := resp.Response.(map[string]interface{})["upload_url"].(string)

	return resp, uploadUrl, err
}

func UploadVideo(filePath, uploadURL string) (map[string]interface{}, error) {

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("video", file.Name())
	if err != nil {
		fmt.Println("Ошибка создания части multipart:", err)
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Ошибка копирования файла:", err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		fmt.Println("Ошибка закрытия multipart:", err)
		return nil, err
	}

	request, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return nil, err
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return nil, err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	return result, nil
}

func (c *VkClient) GetAttachments(filePathsPhoto, filePathsVideo []string, params map[string]string) (string, error) {
	result := ""
	resp, uploadUrlVideo, err := c.VideoSave(params["ownerToken"])
	if err != nil {
		fmt.Println(resp)
		return "", nil
	}
	for _, path := range filePathsVideo {
		resp, err := UploadVideo(path, uploadUrlVideo)
		if err != nil {
			return "", nil
		}

		ownerId := strconv.FormatInt(int64(resp["owner_id"].(float64)), 10)
		ID := strconv.FormatInt(int64(resp["video_id"].(float64)), 10)
		att := fmt.Sprintf("video%s_%s,", ownerId, ID)
		result += att
	}
	resp, uploadUrlPhoto, err := c.GetWallUploadServer(params["groupId"], params["ownerToken"])
	if err != nil {
		fmt.Println(resp)
		return "", nil
	}
	for _, path := range filePathsPhoto {
		resps, err := UploadPhoto(path, uploadUrlPhoto)
		if err != nil {
			return "", nil
		}

		resp, err = c.PhotosSaveWallPhoto(resps, params["ownerToken"], params["groupId"])
		if err != nil {
			return "", nil
		}

		ownerId := strconv.FormatInt(int64(resp.Response.([]interface{})[0].(map[string]interface{})["owner_id"].(float64)), 10)
		ID := strconv.FormatInt(int64(resp.Response.([]interface{})[0].(map[string]interface{})["id"].(float64)), 10)
		att := fmt.Sprintf("photo%s_%s,", ownerId, ID)
		result += att
	}
	result = result[0 : len(result)-1]

	return result, nil
}

func (c *VkClient) WallCloseComments(params map[string]string, token string) (VkResponse, error) {
	return c.CallMethod("wall.closeComments", params, token)
}

func (c *VkClient) PostWithOpt(filePathsPhoto, filePathsVideo []string, params map[string]string, postFromGroup bool, commentFromGroup bool, addEditor bool, closeComment bool, groupJoin bool, carousel bool) (VkResponse, error) {

	//Вступает в группу
	if groupJoin {
		paramsGroupJoin := map[string]string{
			"group_id": params["groupId"],
		}
		resp, err := c.GroupsJoin(paramsGroupJoin, params["token"])
		if err != nil {
			fmt.Println(resp)
			return VkResponse{}, err
		}
	}
	//Устанавливает пользователя редактором
	if addEditor {
		paramsEditManager := map[string]string{
			"group_id": params["groupId"],
			"user_id":  params["userID"],
			"role":     "editor",
		}

		resp, err := c.GroupsEditManager(paramsEditManager, params["token"])
		if err != nil {
			return resp, err
		}
	}

	//Оставляет запись на стене сообщества
	att, err := c.GetAttachments(filePathsPhoto, filePathsVideo, params)
	if err != nil {
		return VkResponse{}, err
	}

	paramsWallPost := map[string]string{
		"owner_id":    params["clientID"],
		"from_group":  "0",
		"attachments": att,
	}
	if postFromGroup {
		paramsWallPost["from_group"] = params["groupId"]
	}
	if carousel {
		paramsWallPost["primary_attachments_mode"] = "carousel"
	}

	resp, postID, link, err := c.WallPost(paramsWallPost, params["messageText"], params["ownerToken"])
	if err != nil {
		return resp, err
	}

	//пишет комментарий под постом от имени сообщества
	paramsCreateComment := map[string]string{
		"owner_id":   params["clientID"],
		"post_id":    postID,
		"from_group": "0",
	}
	if commentFromGroup {
		paramsCreateComment["from_group"] = params["groupId"]
	}
	resp, err = c.WallCreateComment(paramsCreateComment, params["commentText"], params["token"])
	if err != nil {
		return resp, err
	}

	if closeComment {
		//Закрываем комментарии под постом
		paramsCloseComments := map[string]string{
			"owner_id": params["clientID"],
			"post_id":  postID,
		}

		resp, err = c.WallCloseComments(paramsCloseComments, params["ownerToken"])
		if err != nil {
			return resp, err
		}
	}

	//удаляет все роли у выбранного пользователя
	if addEditor {
		paramsEditManager := map[string]string{
			"group_id": params["groupId"],
			"user_id":  params["userID"],
			"role":     "",
		}

		resp, err := c.GroupsEditManager(paramsEditManager, params["ownerToken"])
		if err != nil {
			return resp, err
		}
	}
	//входит из группы
	if groupJoin {
		paramsGroupLeave := map[string]string{
			"group_id": params["groupId"],
		}
		resp, err := c.GroupsLeave(paramsGroupLeave, params["token"])
		if err != nil {
			fmt.Println(resp)
			return VkResponse{}, err
		}
	}
	resp.PostLink = link
	resp.PostID = postID
	return resp, nil
}

func (c *VkClient) TimerPostByMe(targetTime time.Time, params map[string]string, imagePaths, videoPaths []string) (VkResponse, error) {

	targetTime = targetTime.Add(-3 * time.Hour)

	// Вычисляем разницу между текущим моментом и целевой датой
	timeToWait := targetTime.Sub(time.Now())

	timer := time.NewTimer(timeToWait)
	fmt.Println("targetTime: ", targetTime)
	fmt.Println("timeToWait: ", timeToWait)
	fmt.Println("timer: ", timer)
	// Создаем каналы для возврата результата и ошибки
	resultChan := make(chan VkResponse)
	errorChan := make(chan error)

	go func() {
		<-timer.C
		fmt.Println("Таймер сработал")
		resp, err := c.PostWithOptByMe(imagePaths, videoPaths, params, true, false)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- resp
		}
		close(resultChan)
		close(errorChan)
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return VkResponse{}, err
	}
}

func (c *VkClient) GroupsJoin(params map[string]string, token string) (VkResponse, error) {
	return c.CallMethod("groups.join", params, token)
}

func (c *VkClient) GroupsLeave(params map[string]string, token string) (VkResponse, error) {
	return c.CallMethod("groups.leave", params, token)
}

func (c *VkClient) TimerPost(targetTime time.Time, params map[string]string, imagePaths, videoPaths []string) (VkResponse, error) {

	targetTime = targetTime.Add(-3 * time.Hour)

	timeToWait := targetTime.Sub(time.Now())

	timer := time.NewTimer(timeToWait)
	fmt.Println("targetTime: ", targetTime)
	fmt.Println("timeToWait: ", timeToWait)
	fmt.Println("timer: ", timer)

	resultChan := make(chan VkResponse)
	errorChan := make(chan error)

	go func() {
		<-timer.C
		fmt.Println("Таймер сработал")
		resp, err := c.PostWithOpt(imagePaths, videoPaths, params, true, false, false, true, false, false)
		if err != nil {
			errorChan <- err
		} else {
			resultChan <- resp
		}
		close(resultChan)
		close(errorChan)
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errorChan:
		return VkResponse{}, err
	}
}

func (c *VkClient) PostWithOptByMe(filePathsPhoto, filePathsVideo []string, params map[string]string, closeComment bool, carousel bool) (VkResponse, error) {

	//Оставляет запись на стене сообщества
	att, err := c.GetAttachments(filePathsPhoto, filePathsVideo, params)
	if err != nil {
		return VkResponse{}, err
	}

	paramsWallPost := map[string]string{
		"owner_id":    params["clientID"],
		"from_group":  params["groupId"],
		"attachments": att,
	}

	if carousel {
		paramsWallPost["primary_attachments_mode"] = "carousel"
	}

	resp, postID, link, err := c.WallPost(paramsWallPost, params["messageText"], params["ownerToken"])
	if err != nil {
		return resp, err
	}

	//пишет комментарий под постом от имени сообщества
	paramsCreateComment := map[string]string{
		"owner_id":   params["clientID"],
		"post_id":    postID,
		"from_group": "0",
	}

	resp, err = c.WallCreateComment(paramsCreateComment, params["commentText"], params["token"])
	if err != nil {
		return resp, err
	}

	if closeComment {
		//Закрываем комментарии под постом
		paramsCloseComments := map[string]string{
			"owner_id": params["clientID"],
			"post_id":  postID,
		}

		fmt.Println("paramsCloseComments: ", paramsCloseComments)
		resp, err = c.WallCloseComments(paramsCloseComments, params["ownerToken"])
		if err != nil {
			return resp, err
		}
	}
	resp.PostID = postID
	resp.PostLink = link

	return resp, nil
}
