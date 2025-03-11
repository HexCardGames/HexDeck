package api

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SPAMiddleware(fs embed.FS, prefix string, notFoundPath string) gin.HandlerFunc {
	fileServer := http.FileServerFS(fs)

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Next()
			return
		}
		c.Request.URL.Path = prefix + c.Request.URL.Path
		_, err := fs.Open(c.Request.URL.Path)
		if err != nil {
			c.Request.URL.Path = prefix + notFoundPath
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
