package youtube

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/net/html"
)

func DownloadPlaylist() {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println("download to dir=", currentDir)
	set := make(map[string]bool)
	c := make(chan string)

	resp, _ := http.Get("https://www.youtube.com/playlist?list=PLDwdWAXmoUpii8stEkuKPguv4SHhWcHQe")
	z := html.NewTokenizer(resp.Body)
	end := false
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			end = true
		case tt == html.StartTagToken:
			// default:
			t := z.Token()
			re := regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`)
			if re.MatchString(t.String()) {
				url := re.FindStringSubmatch(t.String())[0]
				if !set[url] {
					set[url] = true
					fmt.Println(url)
					y := NewYoutube(true)
					go y.DownloadPlaylistVideo(url, currentDir, c)
				} else {
					fmt.Println(url, "already found")
				}

			}

		}
		if end {
			break
		}
	}

	fmt.Println(len(set), "VIDEOS IN PLAYLIST")
	for i := 0; i < len(set); i++ {
		select {
		case s := <-c:
			fmt.Println(s)
		}
	}
}
