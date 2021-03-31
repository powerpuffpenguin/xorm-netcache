package daemon

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"

	format = `2006/01/02 - 15:04:05`
)

// Run run as deamon
func Run(addr, certFile, keyFile string, debug bool) {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	runGRPC(addr, certFile, keyFile)
}
func printLog(ctx context.Context, at time.Time, fullMethod string, e error, stream bool) {
	var addr string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		strs := md.Get(`x-forwarded-for`)
		if strs != nil {
			addr = strings.Join(strs, `,`)
		}
	}
	if addr == `` {
		if pr, ok := peer.FromContext(ctx); ok {
			addr = pr.Addr.String()
		}
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprint(`[GRPC] `, time.Now().Format(format), ` | `,
		time.Since(at), ` | `,
	))
	if e == nil {
		buffer.WriteString(fmt.Sprint(colorGreen, `success`, colorReset))
	} else {
		buffer.WriteString(fmt.Sprint(colorRed, e, colorReset))
	}
	buffer.WriteString(fmt.Sprint(` | `, addr, ` | `))
	if stream {
		buffer.WriteString(fmt.Sprint(colorYellow, `stream`, colorReset))
	} else {
		buffer.WriteString(fmt.Sprint(colorYellow, `unary`, colorReset))
	}
	buffer.WriteString(fmt.Sprintln(" ->", fullMethod))
	fmt.Print(buffer.String())
}
func unaryLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (response interface{}, e error) {
	at := time.Now()
	response, e = handler(ctx, req)
	printLog(ctx, at, info.FullMethod, e, false)
	return
}
func streamLog(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (e error) {
	at := time.Now()
	ctx := ss.Context()
	e = handler(srv, ss)
	printLog(ctx, at, info.FullMethod, e, true)
	return
}
