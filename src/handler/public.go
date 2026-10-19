package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aaron-vasilev/diary/public"
	"github.com/aaron-vasilev/diary/src/utils"
	"github.com/labstack/echo/v4"
)

type publicFile struct {
	data        []byte
	etag        string
	contentType string
}

// Embedded files never change while the process runs, so prod caches them.
var publicCache sync.Map

// In dev read from disk so tailwind --watch output is picked up without a rebuild.
func publicFS() fs.FS {
	if utils.IsProd() {
		return public.FS
	}
	return os.DirFS("public")
}

func loadPublicFile(name string) (*publicFile, error) {
	if f, ok := publicCache.Load(name); ok {
		return f.(*publicFile), nil
	}

	data, err := fs.ReadFile(publicFS(), name)
	if err != nil {
		return nil, err
	}

	sum := sha256.Sum256(data)
	contentType := mime.TypeByExtension(path.Ext(name))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	f := &publicFile{
		data:        data,
		etag:        `"` + hex.EncodeToString(sum[:8]) + `"`,
		contentType: contentType,
	}

	if utils.IsProd() {
		publicCache.Store(name, f)
	}

	return f, nil
}

func (h HandlerCtx) Public(c echo.Context) error {
	name := path.Clean(c.Param("*"))

	if !fs.ValidPath(name) || name == "." || path.Ext(name) == ".go" {
		return c.NoContent(http.StatusNotFound)
	}

	f, err := loadPublicFile(name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	header := c.Response().Header()
	// no-cache: browser may store the file but must revalidate via ETag on every use.
	header.Set("Cache-Control", "no-cache")
	header.Set("ETag", f.etag)
	header.Set(echo.HeaderContentType, f.contentType)

	http.ServeContent(c.Response(), c.Request(), name, time.Time{}, bytes.NewReader(f.data))
	return nil
}
