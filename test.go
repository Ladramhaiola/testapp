package main

//not degrade message from Sasha
import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"time"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// Operator main worker of programm
type Operator struct {
	client  *http.Client
	service *youtube.Service
}

type audio struct {
	path   string
	format string
	body   *youtube.SearchResult
}

const developerKey = "AIzaSyALx7GChiavVgDs_VGNdTcpyU6P6MufRt8"

func main() {
	c := make(chan string)
	operator, err := initOperator(developerKey)
	if err != nil {
		log.Panic(err)
	}

	searchList := []string{"lil peep"}
	for _, item := range searchList {
		go transfer(operator, "C:/Users/andri/Music/", "mp3", item, c)
	}

	for i := 0; i < len(searchList); i++ {
		fmt.Println(<-c)
	}
}

// complete process of downloading
func transfer(op Operator, path, format, target string, c chan string) {
	ic := make(chan string)
	videos := op.search(target, 4)
	printSR(videos)

	for id, item := range videos {
		go yloader(path+item.Snippet.Title, id, format, ic)
	}

	for i := 0; i < len(videos); i++ {
		fmt.Println(<-ic)
	}
	c <- "done in " + time.Now().String()
}

// search for limited amount of results
func (op Operator) search(target string, lim int64) map[string]*youtube.SearchResult {
	call := op.service.Search.List("id, snippet").Q(target).MaxResults(lim)
	resp, err := call.Do()
	if err != nil {
		log.Panic(err)
	}
	videos := make(map[string]*youtube.SearchResult)
	for _, item := range resp.Items {
		if item.Id.Kind == "youtube#video" {
			videos[item.Id.VideoId] = item
		}
	}
	return videos
}

// create new Operator
func initOperator(devKey string) (Operator, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: devKey},
	}
	service, err := youtube.New(client)
	return Operator{client, service}, err
}

// print search results
func printSR(matches map[string]*youtube.SearchResult) {
	for id, item := range matches {
		fmt.Printf("[%s] %v\n", id, item.Snippet.Title)
	}
}

// URLbyID id -> url
func URLbyID(id string) string {
	return "https://www.youtube.com/watch?v=" + id
}

// load worker
func yloader(path, id, format string, c chan string) {
	destPath := path + "." + format
	url := URLbyID(id)
	var cmd *exec.Cmd
	if format == "mp3" {
		cmd = exec.Command("youtube-dl", "-o", destPath, "--extract-audio", "--audio-format", "mp3", url)
	} else {
		cmd = exec.Command("youtube-dl", "-o", destPath, url)
	}

	err := cmd.Run()
	if err != nil {
		c <- "error " + err.Error()
	} else {
		c <- "done"
	}
}
