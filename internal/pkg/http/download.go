package http

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

func Download(Url string, filename string) error {
	client := &http.Client{
		Timeout: time.Minute * 1,
	}
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return fmt.Errorf("[Download] http.NewRequest error. %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("[Download] Do request error. %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("[Download] io.ReadAll error. %w", err)
	}

	if filename == "" {
		filename, err = ParseFilename(Url)
		if err != nil {
			return fmt.Errorf("[Download] ParseFilename error. %w", err)
		}
	}
	return os.WriteFile(filename, data, 0666)
}

func ParseFilename(Url string) (string, error) {
	parseUrl, err := url.Parse(Url)
	if err != nil {
		return "", fmt.Errorf("[ParseFilename] url.Parse error. %w. url: %s", err, Url)
	}
	return path.Base(parseUrl.Path), nil
}
