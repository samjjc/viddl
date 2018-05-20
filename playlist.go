package youtube

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/net/html"
)

func DownloadPlaylist() {
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println("download to dir=", currentDir)
	set := make(map[string]bool)
	c := make(chan string)
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))

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
					y := NewYoutube(false)
					wg.Add(1)
					bar := p.AddBar(int64(100),
						mpb.PrependDecorators(
							// Display our static name with one space on the right
							decor.StaticName("name", len("name")+1, decor.DidentRight),
							// DwidthSync bit enables same column width synchronization
							decor.Percentage(0, decor.DwidthSync),
						),
						mpb.AppendDecorators(
							// Replace our ETA decorator with "done!", on bar completion event
							decor.OnComplete(decor.ETA(3, 0), "done!", 0, 0),
						),
					)
					go func() {
						defer wg.Done()
						max := 100 * time.Millisecond
						for i := 0; i < 100; i++ {
							time.Sleep(time.Duration(rand.Intn(10)+1) * max / 10)
							bar.Increment()
						}
					}()
					go y.DownloadPlaylistVideo(url, currentDir, c, bar, &wg)
				} else {
					fmt.Println(url, "already found")
				}

			}

		}
		if end {
			break
		}
	}

	wg.Add(len(set))

	fmt.Println(len(set), "VIDEOS IN PLAYLIST")
	for i := 0; i < len(set); i++ {
		select {
		case s := <-c:
			fmt.Println(s)
		}
	}
}
