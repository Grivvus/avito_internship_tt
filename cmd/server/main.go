package main

import (
	"context"
	"log"

	"github.com/Grivvus/reviewers/internal/api"
	"github.com/Grivvus/reviewers/internal/handlers"
	"github.com/Grivvus/reviewers/internal/repository"
	"github.com/Grivvus/reviewers/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	r := gin.Default()

	const connString = "postgres://postgres:hackme@db:5432/postgres"
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal(err)
	}

	// swagger-realted
	r.StaticFile("/openapi.yml", "/api/openapi.yml")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/openapi.yml")))
	// swagger-related end

	userRepository := repository.NewUserRepository(conn)
	teamRepository := repository.NewTeamRepository(conn)
	pullRequestRepository := repository.NewPullRequestRepository(conn)

	userService := service.NewUserservice(userRepository, pullRequestRepository)
	teamService := service.NewTeamService(teamRepository, userRepository)
	pullRequestService := service.NewPullRequestService(pullRequestRepository, userRepository)

	userHandler := handlers.NewUserHandler(userService)
	teamHandler := handlers.NewTeamHandler(teamService)
	pullReqeustHandler := handlers.NewPullRequestHandler(pullRequestService)

	rootHandler := handlers.NewRootHandler(userHandler, teamHandler, pullReqeustHandler)

	api.RegisterHandlers(r, rootHandler)

	r.Run(":8080")
}
