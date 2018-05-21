package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	. "github.com/samjjc/viddl"
	"github.com/vbauerster/mpb"
)

func main() {
	isPlaylist := flag.Bool("p", false, "link of youtube playlist instead of single video")
	flag.Parse()

	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// url := "https://www.youtube.com/watch?v=PyIMjnSfTKc&list=PLDwdWAXmoUpii8stEkuKPguv4SHhWcHQe&t=0s&index=2"
	// url := "https://www.youtube.com/playlist?list=PLDwdWAXmoUpii8stEkuKPguv4SHhWcHQe"
	url := flag.Arg(0)

	if *isPlaylist {

		DownloadPlaylist(currentDir, url)

	} else {

		p := mpb.New()
		bar := CreateLoadingBar(p)
		y := NewYoutube(false, bar)

		y.DownloadSingleVideo(currentDir, url)
	}

	//sleep to finish drawing loading bar
	time.Sleep(time.Millisecond * 200)
}
