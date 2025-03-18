package WorkerClient

import (
	"VkBot/GoogleSheets"
	"VkBot/vkclient"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type WorkerClient struct {
	vkClient           vkclient.VkClient
	googleSheetsClient GoogleSheets.Client
}

func NewWorker(vkClient vkclient.VkClient, googleSheetsClient GoogleSheets.Client) WorkerClient {
	return WorkerClient{
		vkClient:           vkClient,
		googleSheetsClient: googleSheetsClient,
	}
}

func (w *WorkerClient) DownloadImage(destDir string, imageURLs ...string) ([]string, error) {
	var paths []string
	for i, imgURL := range imageURLs {
		filename := strconv.Itoa(i) + ".jpg"

		destPath := filepath.Join(destDir, filename)

		err := downloadFile(imgURL, destPath)
		if err != nil {
			return nil, fmt.Errorf("failed to download image from %s: %w", imgURL, err)
		}

		paths = append(paths, destPath)
	}

	return paths, nil
}

func downloadFile(url string, filepath string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func createTimerToday(postInfo GoogleSheets.PostInfo) *time.Timer {
	targetTime := postInfo.PostDate.Add(-3 * time.Hour)
	timeToWait := targetTime.Sub(time.Now())

	fmt.Println("targetTime: ", targetTime)
	fmt.Println("timeToWait WorkerToday: ", timeToWait)
	return time.NewTimer(timeToWait)
}

func createTimerYesterday(postInfo GoogleSheets.PostInfo) *time.Timer {
	targetTime := postInfo.PostDate.Add(-3 * time.Hour)
	targetTime = targetTime.AddDate(0, 0, 1)
	timeToWait := targetTime.Sub(time.Now())

	fmt.Println("targetTime: ", targetTime)
	fmt.Println("timeToWait WorkerYesterday: ", timeToWait)
	return time.NewTimer(timeToWait)
}

func (w *WorkerClient) WorkerToday(postsInfo []GoogleSheets.PostInfo, wg *sync.WaitGroup) {
	defer wg.Done() // Убедимся, что wg.Done() вызывается при завершении функции

	var accessChannel = make(chan struct{}, 1)

	// Запускаем горутины для каждой задачи
	for _, postInfo := range postsInfo {
		postInfoCopy := postInfo

		// Увеличиваем счётчик WaitGroup перед запуском горутины
		wg.Add(1)
		go func(postInfo GoogleSheets.PostInfo) {
			defer wg.Done() // Убедимся, что wg.Done() вызывается при завершении горутины

			destDir := "sources/post12"
			link := postInfoCopy.TgLink
			commentText := "👇Фотки тут👇\n\n" + link + "\n\n" + link + "\n\n👆Нажимай👆"
			params := map[string]string{
				"commentText": commentText,
				"messageText": postInfoCopy.PostText,
				"userID":      w.vkClient.Config.UserID,            //UserID
				"groupToken":  w.vkClient.Config.GroupTokenCosplay, //GroupTokenCosplay
				"clientID":    w.vkClient.Config.ClientIdCosplay,   //ClientIdCosplay
				"groupId":     w.vkClient.Config.GroupIdCosplay,    //GroupIdCosplay
				"token":       w.vkClient.Config.TokenFake,         //TokenFake
				"ownerToken":  w.vkClient.Config.Token,             //Token
			}

			timer := createTimerToday(postInfoCopy)
			<-timer.C
			accessChannel <- struct{}{}
			fmt.Printf("Таймер сработал для поста №%d от %s\n", postInfoCopy.PostCount, postInfoCopy.ListName)

			photosPath, err := w.DownloadImage(destDir, postInfoCopy.PostPhotoLink1, postInfoCopy.PostPhotoLink2)
			if err != nil {
				log.Printf("failed to download images for post %d by %s: %v", postInfoCopy.PostCount, postInfoCopy.ListName, err)
				<-accessChannel
				return
			}

			var videosPath []string
			resp, err := w.vkClient.TimerPost(postInfoCopy.PostDate, params, photosPath, videosPath)
			if err != nil {
				log.Printf("failed to post timer for post %d by %s: %v", postInfoCopy.PostCount, postInfoCopy.ListName, err)
				<-accessChannel
				return
			}

			err = w.googleSheetsClient.SetPostID(postInfoCopy, resp.PostID)
			if err != nil {
				log.Printf("failed to set post ID for post %d by %s: %v", postInfoCopy.PostCount, postInfoCopy.ListName, err)
				<-accessChannel
				return
			}

			err = w.googleSheetsClient.SetPostLink(postInfoCopy, resp.PostLink)
			if err != nil {
				log.Printf("failed to set post link for post %d by %s: %v", postInfoCopy.PostCount, postInfoCopy.ListName, err)
				<-accessChannel
				return
			}

			fmt.Println(resp)
			<-accessChannel
		}(postInfoCopy)
	}
}

func (w *WorkerClient) WorkerYesterday(postsInfo []GoogleSheets.PostInfo, wg *sync.WaitGroup) {
	defer wg.Done() // Убедимся, что wg.Done() вызывается при завершении функции

	var accessChannel = make(chan struct{}, 1)

	// Запускаем горутины для каждой задачи
	for _, postInfo := range postsInfo {
		postInfoCopy := postInfo

		// Увеличиваем счётчик WaitGroup перед запуском горутины
		wg.Add(1)
		go func(postInfo GoogleSheets.PostInfo) {
			defer wg.Done() // Убедимся, что wg.Done() вызывается при завершении горутины

			timer := createTimerYesterday(postInfoCopy)
			<-timer.C

			accessChannel <- struct{}{}
			fmt.Printf("Таймер сработал для поста №%d от %s\n", postInfoCopy.PostCount, postInfoCopy.ListName)

			id := w.vkClient.Config.ClientIdCosplay + "_" + strconv.Itoa(postInfoCopy.PostID)
			views, err := w.vkClient.GetById(id, w.vkClient.Config.TokenFake)
			if err != nil {
				log.Printf("failed to get views for post %d by %s: %v", postInfoCopy.PostCount, postInfoCopy.ListName, err)
				<-accessChannel
				return
			}

			err = w.googleSheetsClient.SetAverageAudienceReach(postInfoCopy, views)
			if err != nil {
				log.Printf("failed to set average audience reach for post %d by %s: %v", postInfoCopy.PostCount, postInfoCopy.ListName, err)
				<-accessChannel
				return
			}

			<-accessChannel
		}(postInfoCopy)
	}
}

func (w *WorkerClient) StartWorker() error {
	var wg sync.WaitGroup

	// Получаем задачи из Google Sheets
	todayPostQueue, yesterdayPostQueue, err := w.googleSheetsClient.CheckSheet()
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Запускаем WorkerToday
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.WorkerToday(todayPostQueue, &wg)
	}()

	// Запускаем WorkerYesterday
	wg.Add(1)
	go func() {
		defer wg.Done()
		w.WorkerYesterday(yesterdayPostQueue, &wg)
	}()

	// Ждём завершения всех горутин
	wg.Wait()
	return nil
}
