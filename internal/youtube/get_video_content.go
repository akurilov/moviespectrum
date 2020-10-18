package youtube

import (
	"errors"
	"fmt"
	yt "github.com/kkdai/youtube/v2"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	client = &yt.Client{}
	log    = logrus.WithFields(logrus.Fields{})
)

func GetVideoContent(videoId string) (*io.ReadCloser, error) {
	var in *io.ReadCloser = nil
	var err error = nil
	videoMetaData, err := client.GetVideo(videoId)
	var resp *http.Response = nil
	if err == nil {
		filteredFormats := filterLowerResolution(filterRgbFormats(videoMetaData.Formats))
		if len(filteredFormats) < 1 {
			err = errors.New("no allowed video format found")
		} else {
			format := filteredFormats[0]
			filteredFormats = nil
			resp, err = client.GetStream(videoMetaData, &format)
			log.Infof("Youtube video %s: selected input format: %s", videoId, format.MimeType)
		}
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

var (
	allowedMimePrefixes = []string{
		"video/mov",
		"video/ogv",
		"video/webm",
	}
)

func filterRgbFormats(formats yt.FormatList) []yt.Format {
	var filteredFormats = make([]yt.Format, 0)
	for _, format := range formats {
		for _, allowedMimePrefix := range allowedMimePrefixes {
			if strings.HasPrefix(format.MimeType, allowedMimePrefix) {
				log.Debugf("format mime type \"%s\" is allowed", format.MimeType)
				filteredFormats = append(filteredFormats, format)
				break
			}
		}
	}
	return filteredFormats
}

func filterLowerResolution(formats []yt.Format) []yt.Format {
	var filteredFormats = make([]yt.Format, 0)
	for _, format := range formats {
		if format.Width <= 360 {
			log.Debugf("video width %d is allowed", format.Width)
			filteredFormats = append(filteredFormats, format)
		}
	}
	return filteredFormats
}
