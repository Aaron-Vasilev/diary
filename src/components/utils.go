package components

import (
	"os"
	"strings"
)

func Url(path string) string {
	url := os.Getenv("BASE_URL")

	if strings.HasSuffix(url, "/") {
		if strings.HasPrefix(path, "/") {
			url += path[1:]
		} else {
			url += path
		}
	} else if strings.HasPrefix(path, "/") {
		url += path
	} else {
		url += "/" + path
	}

	return url
}
