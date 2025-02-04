package vkclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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

func (c *VkClient) WallPost(params map[string]string, message string, token string) (VkResponse, string, error) {
	escapedMessage := url.QueryEscape(message)
	params["message"] = escapedMessage

	resp, err := c.CallMethod("wall.post", params, token)
	if err != nil {
		return VkResponse{}, "", err
	}
	// Приведение значения post_id к типу float64
	postIdFloat := resp.Response.(map[string]interface{})["post_id"].(float64)

	// Преобразование float64 в целое число и затем в строку
	postId := strconv.FormatInt(int64(postIdFloat), 10)

	return resp, postId, nil
}

func (c *VkClient) WallCreateComment(params map[string]string, message string, token string) (VkResponse, error) {
	escapedMessage := url.QueryEscape(message)
	params["message"] = escapedMessage
	return c.CallMethod("wall.createComment", params, token)
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

	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	// Создаем буфер для multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в multipart/form-data
	part, err := writer.CreateFormFile("photo", file.Name())
	if err != nil {
		fmt.Println("Ошибка создания части multipart:", err)
		return nil, err
	}

	// Копируем содержимое файла в part
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Ошибка копирования файла:", err)
		return nil, err
	}

	// Завершаем создание multipart/form-data
	err = writer.Close()
	if err != nil {
		fmt.Println("Ошибка закрытия multipart:", err)
		return nil, err
	}

	// Формируем HTTP-запрос
	request, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return nil, err
	}

	// Устанавливаем заголовок Content-Type
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Чтение ответа
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return nil, err
	}

	// Парсим ответ в карту
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

	// Открываем файл
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return nil, err
	}
	defer file.Close()

	// Создаем буфер для multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Добавляем файл в multipart/form-data
	part, err := writer.CreateFormFile("video", file.Name())
	if err != nil {
		fmt.Println("Ошибка создания части multipart:", err)
		return nil, err
	}

	// Копируем содержимое файла в part
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Ошибка копирования файла:", err)
		return nil, err
	}

	// Завершаем создание multipart/form-data
	err = writer.Close()
	if err != nil {
		fmt.Println("Ошибка закрытия multipart:", err)
		return nil, err
	}

	// Формируем HTTP-запрос
	request, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		fmt.Println("Ошибка создания запроса:", err)
		return nil, err
	}

	// Устанавливаем заголовок Content-Type
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполняем запрос
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка выполнения запроса:", err)
		return nil, err
	}
	defer response.Body.Close()

	// Чтение ответа
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Ошибка чтения ответа:", err)
		return nil, err
	}

	// Парсим ответ в карту
	var result map[string]interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("ошибка разбора ответа: %w", err)
	}

	return result, nil
}

func (c *VkClient) GetAttachments(filePathsPhoto, filePathsVideo []string, params map[string]string) (string, error) {
	result := ""
	resp, uploadUrlVideo, err := c.VideoSave(params["token"])
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
	resp, uploadUrlPhoto, err := c.GetWallUploadServer(params["groupId"], params["token"])
	if err != nil {
		fmt.Println(resp)
		return "", nil
	}
	for _, path := range filePathsPhoto {
		resps, err := UploadPhoto(path, uploadUrlPhoto)
		if err != nil {
			return "", nil
		}

		resp, err = c.PhotosSaveWallPhoto(resps, params["token"], params["groupId"])
		if err != nil {
			return "", nil
		}

		ownerId := strconv.FormatInt(int64(resp.Response.([]interface{})[0].(map[string]interface{})["owner_id"].(float64)), 10)
		ID := strconv.FormatInt(int64(resp.Response.([]interface{})[0].(map[string]interface{})["id"].(float64)), 10)
		att := fmt.Sprintf("photo%s_%s,", ownerId, ID)
		result += att
	}
	result = result[0 : len(result)-1]
	fmt.Println(result)
	return result, nil
}

func (c *VkClient) WallCloseComments(params map[string]string, token string) (VkResponse, error) {
	return c.CallMethod("wall.closeComments", params, token)
}

func (c *VkClient) PostWithOpt(filePathsPhoto, filePathsVideo []string, params map[string]string, postFromGroup bool, commentFromGroup bool, addEditor bool, closeComment bool, groupJoin bool, carousel bool) (VkResponse, error) {
	fmt.Println(params)
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

		resp, err := c.GroupsEditManager(paramsEditManager, params["ownerToken"])
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

	resp, postID, err := c.WallPost(paramsWallPost, params["messageText"], params["token"])
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

		fmt.Println("paramsCloseComments: ", paramsCloseComments)
		resp, err = c.WallCloseComments(paramsCloseComments, params["token"])
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
	return resp, nil
}

func (c *VkClient) TimerPost(targetTime time.Time, params map[string]string, imagePaths, videoPaths []string) (VkResponse, error) {

	targetTime = targetTime.Add(-3 * time.Hour)

	// Вычисляем разницу между текущим моментом и целевой датой
	timeToWait := targetTime.Sub(time.Now())

	timer := time.NewTimer(timeToWait)
	fmt.Println("targetTime: ", targetTime)
	fmt.Println("timeToWait: ", timeToWait)
	fmt.Println("timer: ", timer)
	fmt.Println("Не в го рутине")
	// Создаем каналы для возврата результата и ошибки
	resultChan := make(chan VkResponse)
	errorChan := make(chan error)

	go func() {
		fmt.Println("В го рутине")
		<-timer.C
		fmt.Println("Таймер сработал")
		resp, err := c.PostWithOpt(imagePaths, videoPaths, params, true, true, true, true, true, true)
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
