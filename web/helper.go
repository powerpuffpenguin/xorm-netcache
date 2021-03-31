package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Offered accept Offered
var Offered = []string{
	binding.MIMEJSON,
	binding.MIMEXML,
	binding.MIMEYAML,
}

// Helper impl
type Helper struct {
}

// NegotiateData .
func (h Helper) NegotiateData(c *gin.Context, code int, data interface{}) {
	switch c.NegotiateFormat(Offered...) {
	case binding.MIMEXML:
		c.XML(code, data)
	case binding.MIMEYAML:
		c.YAML(code, data)
	default:
		// default use json
		c.JSON(code, data)
	}
}

// BindURI .
func (h Helper) BindURI(c *gin.Context, obj interface{}) (e error) {
	e = c.ShouldBindUri(obj)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	return
}

// NegotiateError .
func (h Helper) NegotiateError(c *gin.Context, code int, e error) {
	c.String(code, e.Error())
}

// NegotiateErrorString .
func (h Helper) NegotiateErrorString(c *gin.Context, code int, e string) {
	c.String(code, e)
}

// Bind .
func (h Helper) Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return h.BindWith(c, obj, b)
}

// BindWith .
func (h Helper) BindWith(c *gin.Context, obj interface{}, b binding.Binding) (e error) {
	e = c.ShouldBindWith(obj, b)
	if e != nil {
		h.NegotiateError(c, http.StatusBadRequest, e)
		return
	}
	return
}

// CheckWebsocket .
func (h Helper) CheckWebsocket(c *gin.Context) {
	if !c.IsWebsocket() {
		c.Abort()
		return
	}
}

// BindQuery .
func (h Helper) BindQuery(c *gin.Context, obj interface{}) error {
	return h.BindWith(c, obj, binding.Query)
}

// Upgrade .
func (h Helper) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, responseHeader)
}
