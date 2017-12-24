package main

//not degrade message from Sasha
import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

var (
	query      = enterRequest("lil peep")
	maxResults = setMaxResult(2)
)

// Video implementation
type Video struct {
	Title     string
	url       string
	quality   string
	extension string
}

const developerKey = "AIzaSyALx7GChiavVgDs_VGNdTcpyU6P6MufRt8"

func main() {
	flag.Parse()
	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating new YouTube search: %v", err)
	}

	// Make the API call to YouTube
	call := service.Search.List("id, snippet").Q(*query).MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error making search API call: %v", err)
	}

	// Group video, channel and playlist results to separate lists.
	videos := make(map[string]string)
	channels := make(map[string]string)
	playlist := make(map[string]string)

	// Iterate through each item and add it to the correct list.
	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			videos[item.Id.VideoId] = item.Snippet.Title
		case "youtube#channel":
			channels[item.Id.ChannelId] = item.Snippet.Title
		case "youtube#playlist":
			playlist[item.Id.PlaylistId] = item.Snippet.Title
		}
	}

	prinIDs("Videos", videos)
	prinIDs("Channels", channels)
	prinIDs("Playlists", playlist)
	download("C0DPdy98e4c", "C:\\selflearning\\lol.mp4")
}

// Print ID and title of each result in a list as well as a name that
// identifies the list. For example, print the word section "Videos"
// above a list of video search results, followed by the video ID and title
// of each matching video.
func prinIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("Error occured: %v", err)
	}
}

func enterRequest(req string) *string {
	return flag.String("query", req, "Search term")
}

func setMaxResult(lim int64) *int64 {
	return flag.Int64("max-results", lim, "Max YouTube results")
}

// GetHTTPFromURL initialize a GET request
func GetHTTPFromURL(url string) (body io.ReadCloser, err error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body = response.Body
	if err != nil {
		return nil, err
	}
	return body, err
}

// URLbyID id -> url
func URLbyID(id string) string {
	return "https://youtu.be/" + id
}

func download(id, fname string) {
	url := URLbyID(id)
	body, err := GetHTTPFromURL(url)
	defer body.Close()
	handleError(err)

	output, err := os.Create(fname)
	handleError(err)

	n, err := io.Copy(output, body)
	handleError(err)

	fmt.Println(n, "bytes downloaded")
}

//https://youtu.be/WvV5TbJc9tQ
