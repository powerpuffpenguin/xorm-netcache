package v1

import (
	"fmt"
	"runtime"
	"time"

	"github.com/powerpuffpenguin/xormcache/version"
	"github.com/powerpuffpenguin/xormcache/web"

	"github.com/gin-gonic/gin"
)

var startAt time.Time

// Other .
type Other struct {
	web.Helper
}

// Register impl IHelper
func (h Other) Register(router *gin.RouterGroup) {
	startAt = time.Now()
	router.GET(`version`, h.version)
}
func (h Other) version(c *gin.Context) {
	h.NegotiateObject(c, startAt, gin.H{
		`platform`: fmt.Sprintf(`%s %s %s gin-%s`, runtime.GOOS, runtime.GOARCH, runtime.Version(), gin.Version),
		`tag`:      version.Tag,
		`commit`:   version.Commit,
		`date`:     version.Date,
		`startAt`:  startAt.Unix(),
	})
}
