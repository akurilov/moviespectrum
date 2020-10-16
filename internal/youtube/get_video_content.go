package youtube

import (
	"fmt"
	yt "github.com/kkdai/youtube/v2"
	"io"
	"net/http"
)

func GetVideoContent(client *yt.Client, videoId string) (*io.ReadCloser, error) {
	var in *io.ReadCloser = nil
	var err error = nil
	videoMetaData, err := client.GetVideo(videoId)
	var resp *http.Response = nil
	if err == nil {
		format := &videoMetaData.Formats[0]
		resp, err = client.GetStream(videoMetaData, format)
	}
	if err == nil {
		if resp == nil {
			err = fmt.Errorf("no response, video id: %s", videoId)
		} else {
			in = &resp.Body
		}
	}
	return in, err
}
