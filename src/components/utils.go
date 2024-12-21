package components

import (
	"os"
)

func Url(path string) string {
	return os.Getenv("BASE_URL") + path
}
