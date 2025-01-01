package main

import (
	_ "github.com/goodfoodcesi/tracking-api/docs"
	"github.com/goodfoodcesi/tracking-api/pkg/api"
	"github.com/goodfoodcesi/tracking-api/pkg/config"
)

// @title         Tracking API
// @version       1.0
// @description   Testing Swagger APIs.
// @termsOfService http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /tracking-api
// @schemes   http https
func main() {
	loadConfig := config.LoadConfig()

	//if loadConfig.Env != "dev" {
	//	tracer.Start(
	//		tracer.WithService("tracking-api"),
	//		tracer.WithEnv(loadConfig.Env),
	//		tracer.WithServiceVersion("0.0.5"),
	//	)
	//	defer tracer.Stop()
	//	gin.DefaultWriter = io.Discard
	//}

	r := api.SetupApi(loadConfig)

	if err := r.Run(); err != nil {
		panic(err)
	}
}
