package main

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	httpServer *http.Server
	handler    http.Handler
}

// NewHTTPServer 创建服务器
func NewHTTPServer(config *ServerConfig, handler http.Handler) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)

	return &HTTPServer{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  5 * time.Minute,
			WriteTimeout: 5 * time.Minute,
			IdleTimeout:  120 * time.Second,
		},
		handler: handler,
	}
}

// Start 启动服务器
func (s *HTTPServer) Start() error {
	fmt.Printf("Server starting on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop 停止服务器
func (s *HTTPServer) Stop() error {
	return s.httpServer.Close()
}
