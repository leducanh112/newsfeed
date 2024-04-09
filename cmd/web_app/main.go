package main

import (
	"flag"
	"log"

	"github.com/leducanh112/newsfeed/configs"
	_ "github.com/leducanh112/newsfeed/docs"
	"github.com/leducanh112/newsfeed/internal/app/web_app"
	"github.com/leducanh112/newsfeed/internal/app/web_app/service"
)

// command line flag
var (
	path = flag.String("c", "config.yml", "config path for this service")
)

// @title           Gin Social network Service
// @version         1.0
// @description     A simple social network management service API in Go using Gin framework.
// @termsOfService

// @contact.name   Dong Truong
// @contact.url    https://www.linkedin.com/in/dong-truong-56297a145/
// @contact.email  tpdongcs@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath  /api/v1
//	@securitydefinitions.oauth2.password	OAuth2Password
//	@tokenUrl								https://example.com/oauth/token
//	@scope.read								Grants read access
//	@scope.write							Grants write access
//	@scope.admin							Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode	OAuth2AccessCode
// @tokenUrl								https://example.com/oauth/token
// @authorizationUrl						https://example.com/oauth/authorize
// @scope.admin							Grants read and write access to administrative information
func main() {
	flag.Parse()

	conf, err := configs.GetWebConfig(*path)
	log.Printf("parsed config: %+v", conf)

	if err != nil {
		log.Fatalf("failed to parse config: %s", err)
	}

	webSvc, err := service.NewWebService(conf)
	if err != nil {
		log.Fatalf("failed to init web service: %s", err)
	}

	httpServer := web_app.NewHttpServer(conf.Port, webSvc)

	if err = httpServer.Serve(); err != nil {
		log.Fatalf("failed to run http server: %s", err)
	}
}
