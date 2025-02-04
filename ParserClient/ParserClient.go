package parserclient

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
)

// Interface описывает методы, которые должен реализовать ParserClient
type Interface interface {
	ReadTextFile(filename string) ([]byte, error)
	ParseTexts(data []byte) []string
	ParseLinks(data []byte) []string
	ParseTokens(data []byte) []string
	ParseTimeList(data []byte) []time.Time
	ProcessDirectory(dirPath string) (Dir, error)
	ParsSources() (Sources, error)
}

// ParserClient реализует Interface
type ParserClient struct {
	Sources
}

// Sources представляет собой структуру для хранения информации о директориях
type Sources struct {
	Dirs []Dir
}

// Dir представляет собой структуру для хранения информации о конкретной директории
type Dir struct {
	DirName  string
	TimeList []time.Time
	Texts    []string
	Links    []string
	Tokens   []Token
	Photos   []string
	Videos   []string
}

type Token struct {
	UserTokens []string
	GroupToken string
	OwnerToken string
}

// NewParserClient создает новый экземпляр ParserClient
func NewParserClient() *ParserClient {
	return &ParserClient{Sources{}}
}

// ReadTextFile читает содержимое текстового файла
func (pc *ParserClient) ReadTextFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ParseTexts парсит содержимое text.txt
func (pc *ParserClient) ParseTexts(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	return lines
}

// ParseLinks парсит содержимое links.txt
func (pc *ParserClient) ParseLinks(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	return lines
}

// ParseTokens парсит содержимое tokens.txt
func (pc *ParserClient) ParseTokens(data []byte) []Token {
	var tokens []Token

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			token := Token{
				GroupToken: fields[0],
				OwnerToken: fields[1],
				UserTokens: []string{
					fields[2],
				},
			}
			fmt.Println(fields[0])
			fmt.Println(fields[1])
			fmt.Println(fields[2])
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// ParseTimeList парсит содержимое timeList.txt
func (pc *ParserClient) ParseTimeList(data []byte) []time.Time {
	// Здесь должна быть реализация парсинга времен
	// В качестве примера просто вернем пустой срез
	return []time.Time{}
}

// ProcessDirectory обрабатывает одну директорию
func (pc *ParserClient) ProcessDirectory(dirPath string) (Dir, error) {
	dir := Dir{
		DirName: filepath.Base(dirPath),
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return Dir{}, err
	}

	for _, file := range files {
		filename := file.Name()
		filepath := filepath.Join(dirPath, filename)

		switch {
		case strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".png"):
			dir.Photos = append(dir.Photos, "./sources/"+dir.DirName+"/"+filename)
		case strings.HasSuffix(filename, ".mp4"):
			dir.Videos = append(dir.Videos, "./sources/"+dir.DirName+"/"+filename)
		case filename == "texts.txt":
			data, err := pc.ReadTextFile(filepath)
			if err != nil {
				return Dir{}, err
			}
			dir.Texts = pc.ParseTexts(data)
		case filename == "links.txt":
			data, err := pc.ReadTextFile(filepath)
			if err != nil {
				return Dir{}, err
			}
			dir.Links = pc.ParseLinks(data)
		case filename == "Tokens.txt":
			data, err := pc.ReadTextFile(filepath)
			if err != nil {
				return Dir{}, err
			}
			dir.Tokens = pc.ParseTokens(data)
		case filename == "timeList.txt":
			data, err := pc.ReadTextFile(filepath)
			if err != nil {
				return Dir{}, err
			}
			dir.TimeList = pc.ParseTimeList(data)
		default:
			continue
		}
	}

	return dir, nil
}

// ParsSources обрабатывает директорию sources
func (pc *ParserClient) ParsSources() (Sources, error) {
	sources := Sources{}

	rootDir := "./sources"

	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		log.Fatalf("Ошибка при чтении директории: %v", err)
		return Sources{}, err
	}

	for _, file := range files {
		if file.IsDir() {
			dirPath := filepath.Join(rootDir, file.Name())
			dir, err := pc.ProcessDirectory(dirPath)
			if err != nil {
				return Sources{}, err
			}
			sources.Dirs = append(sources.Dirs, dir)
		}
	}

	return sources, nil
}
