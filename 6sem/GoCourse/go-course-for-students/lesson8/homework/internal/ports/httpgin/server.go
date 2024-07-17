package httpgin

import (
	"homework8/internal/app"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
	app  *gin.Engine
}

func CustomMW(c *gin.Context) {
	t := time.Now()

	c.Next()

	latency := time.Since(t)
	status := c.Writer.Status()

	log.Println("latency", latency, "method", c.Request.Method, "path", c.Request.URL.Path, "status", status)
}

func NewHTTPServer(port string, a app.App) Server {
	gin.SetMode(gin.ReleaseMode)
	s := Server{port: port, app: gin.New()}

	s.app.Use(gin.RecoveryWithWriter(&myWriter{writer: os.Stderr}))

	s.app.Use(CustomMW)

	api := s.app.Group("/api/v1")
	AppRouter(api, a)

	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}

type myWriter struct {
	writer io.Writer
}

func (w *myWriter) Write(message []byte) (int, error) {
	_, err := w.writer.Write([]byte("Attention! Panic!\n"))
	if err != nil {
		return 0, err
	}

	n, err := w.writer.Write(message)
	if err != nil {
		return n, err
	}

	_, err = w.writer.Write([]byte("Attention! Panic!\n"))
	if err != nil {
		return n, err
	}

	return n, nil

}
