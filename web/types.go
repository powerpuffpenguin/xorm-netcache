package web

import (
	"github.com/gin-gonic/gin"
)

// IHelper gin 
type IHelper interface {
	Register(*gin.RouterGroup)
}
