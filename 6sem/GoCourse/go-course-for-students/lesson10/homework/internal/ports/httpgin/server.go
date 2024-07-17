package httpgin

import (
	"homework10/internal/app"
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

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)

	app_ := gin.New()

	app_.Use(gin.RecoveryWithWriter(&myWriter{writer: os.Stderr}))

	app_.Use(CustomMW)

	api := app_.Group("/api/v1")
	AppRouter(api, a)

	return &http.Server{Addr: port, Handler: app_}
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
