package api

import (
	"github.com/powerpuffpenguin/xormcache/web"
	v1 "github.com/powerpuffpenguin/xormcache/web/api/v1"

	"github.com/gin-gonic/gin"
)

// BaseURL request base url
const BaseURL = `api`

// Helper path of /app
type Helper struct {
	web.Helper
}

// Register impl IHelper
func (h Helper) Register(router *gin.RouterGroup) {
	r := router.Group(BaseURL)
	ms := []web.IHelper{
		v1.Helper{},
	}
	for _, m := range ms {
		m.Register(r)
	}
}
