package daemon

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"strings"

	"github.com/powerpuffpenguin/xormcache/protocol/cacher"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server 定義服務器
type Server struct {
	http2Server *http2.Server
	httpServer  *http.Server
	grpcServer  *grpc.Server
	httpService *httpService
}

func runGRPC(addr, certFile, keyFile string) {
	// listen tcp
	l, e := net.Listen(`tcp`, addr)
	if e != nil {
		log.Fatalln(e)
	}
	defer l.Close()
	// new server
	var server Server
	go func() {
		ch := make(chan os.Signal, 2)
		signal.Notify(ch,
			os.Interrupt,
			os.Kill,
			syscall.SIGTERM)
		for {
			sig := <-ch
			switch sig {
			case os.Interrupt:
				server.Stop()
				return
			case syscall.SIGTERM:
				server.Stop()
				return
			}
		}
	}()
	// serve
	if certFile != "" && keyFile != "" {
		log.Println(`h2 work`, addr)
		server.ServeTLS(l, certFile, keyFile)
	} else {
		log.Println(`h2c work`, addr)
		server.Serve(l)
	}
}

// Stop Server
func (s *Server) Stop() {
	s.httpServer.Close()
}

// Serve as h2c
func (s *Server) Serve(l net.Listener) error {
	e := s.init(true, l.Addr().String(), "", "")
	if e != nil {
		return e
	}
	s.httpService = newHTTPService()
	s.httpServer.Handler = h2c.NewHandler(s, s.http2Server)
	e = s.httpServer.Serve(l)
	return e
}

// ServeTLS as h2
func (s *Server) ServeTLS(l net.Listener, certFile, keyFile string) error {
	e := s.init(false, l.Addr().String(), certFile, keyFile)
	if e != nil {
		return e
	}
	s.httpService = newHTTPService()
	s.httpServer.Handler = s
	e = s.httpServer.ServeTLS(l, certFile, keyFile)
	return e
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	contextType := r.Header.Get(`Content-Type`)
	if strings.Contains(contextType, `application/grpc`) {
		s.grpcServer.ServeHTTP(w, r) // grpc service
	} else {
		s.httpService.ServeHTTP(w, r) // gin service
	}
}
func (s *Server) init(h2c bool, address, certFile, keyFile string) (e error) {
	var httpServer http.Server
	var http2Server http2.Server
	// configure http2
	e = http2.ConfigureServer(&httpServer, &http2Server)
	if e != nil {
		return
	}
	// new rpc server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(unaryLog),
		grpc.StreamInterceptor(streamLog),
	)
	// register service
	cacher.RegisterCacherServer(grpcServer, Cacher{})
	// register grpc reflection
	reflection.Register(grpcServer)

	s.httpServer = &httpServer
	s.http2Server = &http2Server
	s.grpcServer = grpcServer
	return
}
