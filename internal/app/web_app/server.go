package web_app

import (
	"fmt"
	"log"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/leducanh112/newsfeed/internal/app/web_app/service"
	v1 "github.com/leducanh112/newsfeed/internal/app/web_app/v1"
)

type HttpServer struct {
	WebService *service.WebService
	Port       int
}

func NewHttpServer(port int, webSvc *service.WebService) *HttpServer {
	return &HttpServer{
		WebService: webSvc,
		Port:       port,
	}
}

func (c *HttpServer) Serve() error {
	r := gin.Default()

	v1Router := r.Group("/api/v1")
	v1.AddUserRouter(v1Router, c.WebService)
	v1.AddPostRouter(v1Router, c.WebService)

	initSwagger(r)
	initPprof(r)
	initPrometheus(r)

	address := fmt.Sprintf(":%d", c.Port) // :port
	log.Println("")
	if err := r.Run(address); err != nil {
		return fmt.Errorf("failed to run http server on %s: %s", address, err)
	}
	return nil
}

func initSwagger(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initPprof(r *gin.Engine) {
	r.GET("/debug/pprof/:profile", func(c *gin.Context) {
		pprof.Index(c.Writer, c.Request)
	})
	r.GET("/debug/pprof/profile", func(c *gin.Context) {
		pprof.Profile(c.Writer, c.Request)
	})
	r.GET("/debug/pprof/trace", func(c *gin.Context) {
		pprof.Trace(c.Writer, c.Request)
	})
}

func initPrometheus(r *gin.Engine) {
	handler := promhttp.Handler()
	r.GET("/metrics", func(context *gin.Context) {
		handler.ServeHTTP(context.Writer, context.Request)
	})
}
