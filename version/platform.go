package version

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Platform .
var Platform = runtime.GOOS +
	` ` + runtime.GOARCH +
	` ` + runtime.Version() +
	` grpc` + grpc.Version +
	` gin-` + gin.Version
