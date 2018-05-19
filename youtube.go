package youtube

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type stream map[string]string

type Youtube struct {
	DebugMode         bool
	StreamList        []stream
	VideoID           string
	videoInfo         string
	DownloadPercent   chan int64
	contentLength     float64
	totalWrittenBytes float64
	downloadLevel     float64
}

func (y *Youtube) DecodeURL(url string) error {
	err := y.findVideoID(url)
	if err != nil {
		return fmt.Errorf("findVideoID error=%s", err)
	}

	err = y.getVideoInfo()
	if err != nil {
		return fmt.Errorf("getVideoInfo error=%s", err)
	}

	err = y.parseVideoInfo()
	if err != nil {
		return fmt.Errorf("parse video info failed, err=%s", err)
	}

	return nil
}

func (y *Youtube) parseVideoInfo() error {
	answer, err := url.ParseQuery(y.videoInfo)
	if err != nil {
		return err
	}

	status, ok := answer["status"]
	if !ok {
		err = fmt.Errorf("no response status found in the server's answer")
		return err
	}
	if status[0] == "fail" {
		reason, ok := answer["reason"]
		if ok {
			err = fmt.Errorf("'fail' response status found in the server's answer, reason: '%s'", reason[0])
		} else {
			err = errors.New(fmt.Sprint("'fail' response status found in the server's answer, no reason given"))
		}
		return err
	}
	if status[0] != "ok" {
		err = fmt.Errorf("non-success response status found in the server's answer (status: '%s')", status)
		return err
	}

	// read the streams map
	streamMap, ok := answer["url_encoded_fmt_stream_map"]
	if !ok {
		err = errors.New(fmt.Sprint("no stream map found in the server's answer"))
		return err
	}

	// read each stream
	streamsList := strings.Split(streamMap[0], ",")

	var streams []stream
	for streamPos, streamRaw := range streamsList {
		streamQry, err := url.ParseQuery(streamRaw)
		if err != nil {
			log.Printf("An error occured while decoding one of the video's stream's information: stream %d: %s\n", streamPos, err)
			continue
		}
		var sig string
		if _, exist := streamQry["sig"]; exist {
			sig = streamQry["sig"][0]
		}

		streams = append(streams, stream{
			"quality": streamQry["quality"][0],
			"type":    streamQry["type"][0],
			"url":     streamQry["url"][0],
			"sig":     sig,
			"title":   answer["title"][0],
			"author":  answer["author"][0],
		})
		y.log(fmt.Sprintf("Stream found: quality '%s', format '%s'", streamQry["quality"][0], streamQry["type"][0]))
	}

	y.StreamList = streams
	return nil
}

func (y *Youtube) getVideoInfo() error {
	url := "http://youtube.com/get_video_info?video_id=" + y.VideoID
	y.log(fmt.Sprintf("url: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	y.videoInfo = string(body)
	return nil
}

func (y *Youtube) findVideoID(url string) error {
	videoID := url
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		reList := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range reList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}
	log.Printf("Found video id: '%s'", videoID)
	y.VideoID = videoID
	if strings.ContainsAny(videoID, "?&/<%=") {
		return errors.New("invalid characters in video id")
	}
	if len(videoID) < 10 {
		return errors.New("the video id must be at least 10 characters long")
	}
	return nil
}

func (y *Youtube) log(logText string) {
	if y.DebugMode {
		log.Println(logText)
	}
}
