package infrastructure

import (
	accountsHttp "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/http"
	blocksHttp "github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/http"
	coursesHttp "github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/http"
	laboratoriesHttp "github.com/UPB-Code-Labs/main-api/src/laboratories/infrastructure/http"
	languagesHttp "github.com/UPB-Code-Labs/main-api/src/languages/infrastructure/http"
	rubricsHttp "github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/http"
	sessionHttp "github.com/UPB-Code-Labs/main-api/src/session/infrastructure/http"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var routesGroups = []func(*gin.RouterGroup){
	accountsHttp.StartAccountsRoutes,
	blocksHttp.StartBlocksRoutes,
	sessionHttp.StartSessionRoutes,
	coursesHttp.StartCoursesRoutes,
	rubricsHttp.StartRubricsRoutes,
	laboratoriesHttp.StartLaboratoriesRoutes,
	languagesHttp.StartLanguagesRoutes,
}

func InstanceHttpServer() (r *gin.Engine) {
	engine := gin.Default()
	engine.Use(sharedInfrastructure.ErrorHandlerMiddleware())

	isInProductionEnvironment := sharedInfrastructure.GetEnvironment().Environment == "production"
	if isInProductionEnvironment {
		gin.SetMode(gin.ReleaseMode)
	}

	// Configure CORS rules
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowOrigins = []string{sharedInfrastructure.GetEnvironment().WebClientUrl}
	corsConfig.AllowCredentials = true
	engine.Use(cors.New(corsConfig))

	// Register routes
	baseGroup := engine.Group("/api/v1")
	for _, registerRoutesGroup := range routesGroups {
		registerRoutesGroup(baseGroup)
	}

	return engine
}
