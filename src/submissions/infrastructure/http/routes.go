package http

import (
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/application"
	"github.com/UPB-Code-Labs/main-api/src/submissions/infrastructure/implementations"
	"github.com/gin-gonic/gin"
)

func StartSubmissionsRoutes(g *gin.RouterGroup) {
	submissionsGroup := g.Group("/submissions")

	useCases := application.SubmissionUseCases{
		SubmissionsRepository:   implementations.GetSubmissionsRepositoryInstance(),
		SubmissionsQueueManager: implementations.GetSubmissionsRabbitMQQueueManagerInstance(),
	}

	controllers := SubmissionsController{
		UseCases: &useCases,
	}

	submissionsGroup.POST(
		":test_block_uuid",
		sharedInfrastructure.WithAuthenticationMiddleware(),
		sharedInfrastructure.WithAuthorizationMiddleware([]string{"student"}),
		controllers.HandleReceiveSubmissions,
	)

	/*
		submissionsGroup.GET(
			"/:test_block_uuid/mine",
			sharedInfrastructure.WithAuthenticationMiddleware(),
			sharedInfrastructure.WithAuthorizationMiddleware([]string{"student"}),
			controllers.HandleGetSubmission,
		)
	*/
}
