package daemon

import (
	"net/http"

	"github.com/powerpuffpenguin/xormcache/web"
	"github.com/powerpuffpenguin/xormcache/web/api"

	"github.com/gin-gonic/gin"
)

type httpService struct {
	router *gin.Engine
}

func newHTTPService() *httpService {
	router := gin.Default()
	rs := []web.IHelper{
		api.Helper{},
	}
	for _, r := range rs {
		r.Register(&router.RouterGroup)
	}
	return &httpService{
		router: router,
	}
}

func (s *httpService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
