package components

import (
	"os"
	"strings"
)

func Url(path string) string {
	baseUrl := os.Getenv("BASE_URL")

	if strings.HasSuffix(baseUrl, "/") {
		if strings.HasPrefix(path, "/") {
			return baseUrl + path[1:]
		} else {
			return baseUrl + path
		}
	} else if strings.HasPrefix(path, "/") {
		return baseUrl + path
	} else {
		return baseUrl + "/" + path
	} 
}
