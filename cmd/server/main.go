package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/handlers"
)

func main() {
	r := gin.Default()

	// swagger-realted
	r.StaticFile("/openapi.yml", "./api/openapi.yml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/openapi.yml")))
	// swagger-related end

	userHandler := handlers.NewUserHandler()
	teamHandler := handlers.NewTeamHandler()
	pullReqeustHandler := handlers.NewPullRequestHandler()

	rootHandler := handlers.NewRootHandler(userHandler, teamHandler, pullReqeustHandler)

	api.RegisterHandlers(r, rootHandler)

	r.Run(":8080")
}
