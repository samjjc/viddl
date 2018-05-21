package viddl

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/net/html"
)

func DownloadPlaylist(currentDir string, url string) {
	log.Println("download to dir=", currentDir)
	set := make(map[string]bool)
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	i := 1
	resp, _ := http.Get(url)
	// resp, _ := http.Get("https://www.youtube.com/playlist?list=PLDwdWAXmoUpii8stEkuKPguv4SHhWcHQe")
	// resp, _ := http.Get("https://www.youtube.com/playlist?list=PLDwdWAXmoUpjjkvj6NA5DwX0v5Zg8BpcN")
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
					name := fmt.Sprintf("Bar#%d:", i)
					i++
					set[url] = true
					fmt.Println(url)
					bar := p.AddBar(100,
						mpb.PrependDecorators(
							// Display our static name with one space on the right
							decor.StaticName(name, len(name)+1, decor.DidentRight),
							// DwidthSync bit enables same column width synchronization
							decor.Percentage(0, decor.DwidthSync),
						),
						mpb.AppendDecorators(
							// Replace our ETA decorator with "done!", on bar completion event
							decor.OnComplete(decor.ETA(3, 0), "done!", 0, 0),
						),
					)
					y := NewYoutube(false, bar)
					wg.Add(1)
					go y.DownloadPlaylistVideo(url, currentDir, &wg)
				} else {
					fmt.Println(url, "already found")
				}

			}

		}
		if end {
			break
		}
	}
	wg.Wait()
}
