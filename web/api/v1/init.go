package v1

import (
	"github.com/powerpuffpenguin/xormcache/web"

	"github.com/gin-gonic/gin"
)

// BaseURL .
const BaseURL = `v1`

// Helper .
type Helper struct {
	web.Helper
}

// Register impl IController
func (h Helper) Register(router *gin.RouterGroup) {
	r := router.Group(BaseURL)

	ms := []web.IHelper{
		Other{},
	}
	for _, m := range ms {
		m.Register(r)
	}
}
